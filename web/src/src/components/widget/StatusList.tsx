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
            className="max-h-[300px]"
            data={data.map(({ description, title, running }) => ({
                id: title,
                primary: (
                    <div>
                        <div className="font-medium">{title}</div>
                        <div className="flex items-center space-x-1">
                            <Icon
                                path={mdiToolboxOutline}
                                size={0.6}
                                className="flex-shrink-0 opacity-50"
                            />
                            <span className="text-xs opacity-50">{description}</span>
                        </div>
                    </div>
                ),
                secondary: (
                    <div className={`badge ${running ? 'badge-success' : 'badge-error'}`}>
                        <Icon
                            className="text-white"
                            path={running ? mdiCheckCircle : mdiCloseCircle}
                            size={0.6}
                        />
                    </div>
                )
            }))}
        />
    );
};
