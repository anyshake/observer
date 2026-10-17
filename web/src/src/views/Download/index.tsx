import '@fancyapps/ui/dist/fancybox/fancybox.css';

import { mdiArchive, mdiClose, mdiImageAlbum, mdiMagnify } from '@mdi/js';
import Icon from '@mdi/react';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';

import { Card } from '../../components/widget/Card';
import { TableList } from '../../components/widget/TableList';
import { IRouterComponent } from '../../config/router';
import { useAssetActions } from './useAssetActions';
import { useFileColumns } from './useFileColumns';
import { useFileLists } from './useFileLists';

const Download = ({ currentLocale }: IRouterComponent) => {
    const { t } = useTranslation();
    const { getFileListDataLoading, miniSeedList, helicorderList } = useFileLists();

    const [searchMiniSeed, setSearchMiniSeed] = useState('');
    const [searchHelicorder, setSearchHelicorder] = useState('');
    const filteredMiniSeedList = miniSeedList.filter((item) =>
        item.fileName.toLowerCase().includes(searchMiniSeed.toLowerCase())
    );
    const filteredHelicorderList = helicorderList.filter((item) =>
        item.fileName.toLowerCase().includes(searchHelicorder.toLowerCase())
    );

    const assetActions = useAssetActions();
    const { getColumns } = useFileColumns(assetActions);

    return (
        <div className="app-page">
            <div className="app-column">
                <Card title={t('views.Download.miniseed_list.title')} iconPath={mdiArchive}>
                    {getFileListDataLoading ? (
                        <div className="app-center">
                            <span className="loading loading-spinner loading-lg app-loading-muted" />
                        </div>
                    ) : (
                        <div className="app-stack">
                            <label className="input">
                                <Icon className="app-icon" path={mdiMagnify} size={0.8} />
                                <input
                                    type="search"
                                    className="app-grow"
                                    value={searchMiniSeed}
                                    onChange={({ currentTarget }) =>
                                        setSearchMiniSeed(currentTarget.value)
                                    }
                                    placeholder={t('views.Download.miniseed_list.search')}
                                />
                                {searchMiniSeed.length > 0 && (
                                    <button
                                        className="app-icon-button"
                                        onClick={() => setSearchMiniSeed('')}
                                    >
                                        <Icon className="app-icon" path={mdiClose} size={0.8} />
                                    </button>
                                )}
                            </label>
                            <TableList
                                sortField="id"
                                sortDirection="asc"
                                currentLocale={currentLocale}
                                data={filteredMiniSeedList}
                                columns={getColumns(false)}
                            />
                        </div>
                    )}
                </Card>
                <Card title={t('views.Download.helicorder_list.title')} iconPath={mdiImageAlbum}>
                    {getFileListDataLoading ? (
                        <div className="app-center">
                            <span className="loading loading-spinner loading-lg app-loading-muted" />
                        </div>
                    ) : (
                        <div className="app-stack">
                            <label className="input">
                                <Icon className="app-icon" path={mdiMagnify} size={0.8} />
                                <input
                                    type="search"
                                    className="app-grow"
                                    value={searchHelicorder}
                                    onChange={({ currentTarget }) =>
                                        setSearchHelicorder(currentTarget.value)
                                    }
                                    placeholder={t('views.Download.helicorder_list.search')}
                                />
                                {searchHelicorder.length > 0 && (
                                    <button
                                        className="app-icon-button"
                                        onClick={() => setSearchHelicorder('')}
                                    >
                                        <Icon className="app-icon" path={mdiClose} size={0.8} />
                                    </button>
                                )}
                            </label>
                            <TableList
                                sortField="id"
                                sortDirection="asc"
                                currentLocale={currentLocale}
                                data={filteredHelicorderList}
                                columns={getColumns(true)}
                            />
                        </div>
                    )}
                </Card>
            </div>
        </div>
    );
};

export default Download;
