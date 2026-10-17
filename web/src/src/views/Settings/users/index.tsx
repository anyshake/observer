import { mdiAccountPlus, mdiAccountRemove, mdiClose, mdiMagnify, mdiRefresh } from '@mdi/js';
import Icon from '@mdi/react';
import { GridColDef, GridValidRowModel } from '@mui/x-data-grid';
import { useTranslation } from 'react-i18next';

import { ErrorPage } from '../../../components/ui/ErrorPage';
import { TableList } from '../../../components/widget/TableList';
import { CreateUserModal } from './CreateUserModal';
import { EditUserModal } from './EditUserModal';
import { useUserColumns } from './useUserColumns';
import { useUsers } from './useUsers';

interface IUsers {
    currentLocale: string;
}

export const Users = ({ currentLocale }: IUsers) => {
    const { t } = useTranslation();

    const users = useUsers();
    const {
        getUserListDataError,
        getUserListDataLoading,
        searchUser,
        setSearchUser,
        filteredUserList,
        handleRefreshUserList,
        selectedUsers,
        handleSelectUsers,
        handleRemoveUsers,
        setCreateUserModalOpen
    } = users;
    const { columns } = useUserColumns(users);

    return getUserListDataError ? (
        <ErrorPage
            content={getUserListDataError.message}
            debug={JSON.stringify(getUserListDataError)}
        />
    ) : getUserListDataLoading ? (
        <div className="app-center">
            <span className="loading loading-spinner app-loading-delayed" />
        </div>
    ) : (
        <>
            <div className="app-toolbar">
                <button
                    className="btn btn-sm app-button-toolbar"
                    onClick={() => handleRefreshUserList()}
                >
                    <Icon className="app-icon" path={mdiRefresh} size={0.7} />
                    <span>{t('views.Settings.Users.refresh_list.submit_button')}</span>
                </button>
                <button
                    className="btn btn-sm app-button-toolbar"
                    onClick={() => setCreateUserModalOpen(true)}
                >
                    <Icon className="app-icon" path={mdiAccountPlus} size={0.7} />
                    <span>{t('views.Settings.Users.create_user.submit_button')}</span>
                </button>
                {selectedUsers.length > 1 && (
                    <button
                        className="btn btn-sm app-button-toolbar"
                        onClick={() => handleRemoveUsers()}
                    >
                        <Icon className="app-icon" path={mdiAccountRemove} size={0.7} />
                        <span>{t('views.Settings.Users.remove_users.submit_button')}</span>
                    </button>
                )}
            </div>
            <label className="input">
                <Icon className="app-icon" path={mdiMagnify} size={0.8} />
                <input
                    type="search"
                    className="app-grow"
                    value={searchUser}
                    onChange={({ currentTarget }) => setSearchUser(currentTarget.value)}
                    placeholder={t('views.Settings.Users.search_users.placeholder')}
                />
                {searchUser.length > 0 && (
                    <button className="app-icon-button" onClick={() => setSearchUser('')}>
                        <Icon className="app-icon" path={mdiClose} size={0.8} />
                    </button>
                )}
            </label>

            <TableList
                currentLocale={currentLocale}
                data={filteredUserList}
                columns={columns as GridColDef<GridValidRowModel>[]}
                onSelect={(selected) => handleSelectUsers(selected)}
            />

            <CreateUserModal users={users} />

            <EditUserModal users={users} />
        </>
    );
};
