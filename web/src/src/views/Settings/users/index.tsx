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
        <div className="flex items-center justify-center">
            <span className="loading loading-spinner text-primary mt-6" />
        </div>
    ) : (
        <>
            <div className="flex flex-wrap items-center gap-2">
                <button
                    className="btn btn-sm flex items-center"
                    onClick={() => handleRefreshUserList()}
                >
                    <Icon className="flex-shrink-0" path={mdiRefresh} size={0.7} />
                    <span>{t('views.Settings.Users.refresh_list.submit_button')}</span>
                </button>
                <button
                    className="btn btn-sm flex items-center"
                    onClick={() => setCreateUserModalOpen(true)}
                >
                    <Icon className="flex-shrink-0" path={mdiAccountPlus} size={0.7} />
                    <span>{t('views.Settings.Users.create_user.submit_button')}</span>
                </button>
                {selectedUsers.length > 1 && (
                    <button
                        className="btn btn-sm flex items-center"
                        onClick={() => handleRemoveUsers()}
                    >
                        <Icon className="flex-shrink-0" path={mdiAccountRemove} size={0.7} />
                        <span>{t('views.Settings.Users.remove_users.submit_button')}</span>
                    </button>
                )}
            </div>
            <label className="input">
                <Icon className="flex-shrink-0" path={mdiMagnify} size={0.8} />
                <input
                    type="search"
                    className="grow"
                    value={searchUser}
                    onChange={({ currentTarget }) => setSearchUser(currentTarget.value)}
                    placeholder={t('views.Settings.Users.search_users.placeholder')}
                />
                {searchUser.length > 0 && (
                    <button
                        className="cursor-pointer opacity-60 transition-all hover:opacity-100"
                        onClick={() => setSearchUser('')}
                    >
                        <Icon className="flex-shrink-0" path={mdiClose} size={0.8} />
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
