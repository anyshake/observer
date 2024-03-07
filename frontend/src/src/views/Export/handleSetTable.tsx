import { Dispatch, SetStateAction } from "react";
import { ExportsUpdates } from "./getExportsUpdates";
import { TableProps } from "../../components/Table";
import { getTimeString } from "../../helpers/utils/getTimeString";

export const handleSetTable = (
    res: ExportsUpdates,
    stateFn: Dispatch<SetStateAction<TableProps>>
) => {
    if (!res?.data) {
        stateFn((prev) => ({
            ...prev,
            placeholder: "views.export.table.placeholder.fetch_mseed_error",
        }));
        return;
    }

    const tableData = res.data
        .sort((a, b) => b.time - a.time)
        .map((item) => {
            const timeString = getTimeString(new Date(item.time).getTime());
            return { ...item, time: timeString };
        });
    stateFn((prev) => ({ ...prev, data: tableData }));
};
