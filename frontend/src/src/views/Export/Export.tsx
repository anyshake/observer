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

// Export timeout is 100s by default
const EXPORT_TIMEOUT = 100 * 1000;

interface ExportState {
    table: TableProps;
    progress: ProgressProps;
}

class Export extends Component<WithTranslation, ExportState> {
    constructor(props: WithTranslation) {
        super(props);
        this.state = {
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
            progress: {
                value: 0,
            },
        };
    }

    // Export specified MiniSEED file
    exportMiniSEED = async ({ name }: TableData): Promise<void> => {
        const { t } = this.props;
        await toast.promise(
            restfulApiByTag({
                blob: true,
                tag: "mseed",
                filename: name,
                timeout: EXPORT_TIMEOUT,
                body: { action: "export", name },
                onDownload: (e) => {
                    const { progress } = e;
                    this.setState({
                        progress: { value: (progress || 0) * 100 },
                    });
                },
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

    render() {
        const { table, progress } = this.state;
        const { value } = progress;

        return (
            <View>
                <Header />
                <Sidebar />

                <Content>
                    <Navbar />

                    <Container layout="none">
                        <Card label={{ id: "views.export.cards.file_list" }}>
                            {!!value && <Progress {...progress} />}
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
