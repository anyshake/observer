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
import FloppyIcon from "../../assets/icons/floppy-disk-solid.svg";
import { WithTranslation, withTranslation } from "react-i18next";
import restfulApiByTag from "../../helpers/request/restfulApiByTag";

// Export timeout is 100s by default
const EXPORT_TIMEOUT = 100 * 1000;

interface ExportState {
    table: TableProps;
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
                ],
                placeholder: { id: "views.export.table.placeholder" },
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
            }),
            {
                loading: t("views.export.toasts.is_exporting_mseed"),
                success: t("views.export.toasts.export_mseed_success"),
                error: t("views.export.toasts.export_mseed_error"),
            }
        );
    };

    // Fetch available MiniSEED files
    async componentDidMount(): Promise<void> {
        const { t } = this.props;
        const { table } = this.state;

        // Read MiniSEED file list from server
        const loader = toast.loading(
            t("views.export.toasts.is_fetching_mseed")
        );
        const { data } = await restfulApiByTag({
            tag: "mseed",
            body: { action: "show" },
            timeout: EXPORT_TIMEOUT,
        });
        toast.remove(loader);

        // Check if there is no data
        if (!data) {
            toast.error(t("views.export.toasts.fetch_mseed_error"));
            return;
        }

        // Append export action and show success message
        const actions = [
            {
                icon: FloppyIcon,
                onClick: this.exportMiniSEED,
                label: { id: "views.export.table.actions.export" },
            },
        ];
        this.setState({ table: { ...table, data, actions } });
        toast.success(t("views.export.toasts.fetch_mseed_success"));
    }

    render() {
        const { table } = this.state;
        return (
            <View>
                <Header />
                <Sidebar />

                <Content>
                    <Navbar />

                    <Container layout="none">
                        <Card label={{ id: "views.export.cards.file_list" }}>
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
