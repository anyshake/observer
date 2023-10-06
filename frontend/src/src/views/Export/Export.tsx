import { Component } from "react";
import Content from "../../components/Content";
import Header from "../../components/Header";
import Navbar from "../../components/Navbar";
import Sidebar from "../../components/Sidebar";
import View from "../../components/View";
import Scroller from "../../components/Scroller";
import Footer from "../../components/Footer";
import toast, { Toaster } from "react-hot-toast";
import Card from "../../components/Card";
import Container from "../../components/Container";
import Table, { TableData, TableProps } from "../../components/Table";
import ExportIcon from "../../assets/icons/download-solid.svg";
import { WithTranslation, withTranslation } from "react-i18next";
import restfulApiByTag from "../../helpers/request/restfulApiByTag";
import Progress, { ProgressProps } from "../../components/Progress";
import getSortedArray from "../../helpers/array/getSortedArray";
import axios, { CancelTokenSource } from "axios";

// Export timeout is 100s by default
const EXPORT_TIMEOUT = 100 * 1000;

interface ExportState {
    readonly table: TableProps;
    readonly tasks: ProgressProps[];
    readonly tokens: CancelTokenSource[];
}

class Export extends Component<WithTranslation, ExportState> {
    constructor(props: WithTranslation) {
        super(props);
        this.state = {
            tokens: [],
            tasks: [],
            table: {
                data: [],
                actions: [],
                columns: [
                    {
                        key: "name",
                        label: { id: "views.export.table.columns.name" },
                    },
                    {
                        key: "size",
                        label: { id: "views.export.table.columns.size" },
                    },
                    {
                        key: "time",
                        label: { id: "views.export.table.columns.time" },
                    },
                    {
                        key: "ttl",
                        label: { id: "views.export.table.columns.ttl" },
                    },
                ],
                placeholder: { id: "views.export.table.placeholder" },
            },
        };
    }

    // Update export task progress by name
    updateTaskProgress = (name: string, value: number): void => {
        // Check if task is new
        const { tasks } = this.state;
        const taskIndex = tasks.findIndex(({ label }) => label === name);
        if (taskIndex === -1) {
            // Add new task if not found
            tasks.push({ label: name, value });
        } else if (value === 100) {
            // Remove task if completed
            setTimeout(() => {
                tasks.splice(taskIndex, 1);
                this.setState({ tasks });
            }, 1000);
        } else {
            // Update task progress
            tasks[taskIndex].value = value;
            this.setState({ tasks });
        }
    };

    // Export specified MiniSEED file
    exportMiniSEED = async ({ name }: TableData): Promise<void> => {
        // Create cancel token and add to list
        const { tokens } = this.state;
        const { source } = axios.CancelToken;
        const cancelToken = source();
        tokens.push(cancelToken);

        // Show toast and update task progress
        const { t } = this.props;
        await toast.promise(
            restfulApiByTag({
                cancelToken,
                blob: true,
                tag: "mseed",
                filename: name,
                timeout: EXPORT_TIMEOUT,
                body: { action: "export", name },
                onDownload: ({ progress }) =>
                    this.updateTaskProgress(name, (progress || 0) * 100),
            }),
            {
                loading: t("views.export.toasts.is_exporting_mseed"),
                success: t("views.export.toasts.export_mseed_success"),
                error: t("views.export.toasts.export_mseed_error"),
            }
        );
    };

    // Fetch available MiniSEED files from server
    async componentDidMount(): Promise<void> {
        const { t } = this.props;
        const { table } = this.state;

        // Read MiniSEED file list from server
        const { data } = await restfulApiByTag({
            tag: "mseed",
            body: { action: "show" },
            timeout: EXPORT_TIMEOUT,
        });
        if (!data || !data.length) {
            // Show error toast if no MiniSEED files found
            const err = "views.export.toasts.fetch_mseed_error";
            toast.error(t(err));
            this.setState({
                table: { ...table, placeholder: { id: err } },
            });
        } else {
            // Sort files by time and show success toast
            toast.success(t("views.export.toasts.fetch_mseed_success"));
            getSortedArray(data, "time", "datetime", "desc");
            // Append export action and update table
            const actions = [
                {
                    icon: ExportIcon,
                    onClick: this.exportMiniSEED,
                    label: { id: "views.export.table.actions.export" },
                },
            ];
            this.setState({ table: { ...table, data, actions } });
        }
    }

    // Cancel all pending requests
    componentWillUnmount(): void {
        const { tokens } = this.state;
        tokens.forEach(({ cancel }) => cancel());
    }

    render() {
        const { table, tasks } = this.state;

        return (
            <View>
                <Header />
                <Sidebar />

                <Content>
                    <Navbar />

                    <Container layout="none">
                        <Card label={{ id: "views.export.cards.file_list" }}>
                            {tasks.map(
                                (item, index) =>
                                    !!item.value && (
                                        <Progress key={index} {...item} />
                                    )
                            )}
                            <Table {...table} />
                        </Card>
                    </Container>
                </Content>

                <Scroller />
                <Footer />
                <Toaster position="top-center" />
            </View>
        );
    }
}

export default withTranslation()(Export);
