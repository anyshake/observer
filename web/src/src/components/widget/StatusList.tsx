import { mdiCheckCircle, mdiCloseCircle, mdiToolboxOutline } from '@mdi/js';
import Icon from '@mdi/react';

import { List } from './List';

interface IStatusList {
    readonly data: Array<{
        readonly description: string;
        readonly title: string;
        readonly running: boolean;
    }>;
}

export const StatusList = ({ data }: IStatusList) => {
    return (
        <List
            className="app-status-list"
            data={data.map(({ description, title, running }) => ({
                id: title,
                primary: (
                    <div>
                        <div className="app-text-medium">{title}</div>
                        <div className="app-row-tight">
                            <Icon path={mdiToolboxOutline} size={0.6} className="app-status-icon" />
                            <span className="app-status-description">{description}</span>
                        </div>
                    </div>
                ),
                secondary: (
                    <div
                        className={`badge app-status-badge ${running ? 'badge-success' : 'badge-error'}`}
                    >
                        <Icon
                            className="app-icon"
                            path={running ? mdiCheckCircle : mdiCloseCircle}
                            size={0.6}
                        />
                    </div>
                )
            }))}
        />
    );
};
