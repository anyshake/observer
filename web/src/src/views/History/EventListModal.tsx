import { mdiTarget } from '@mdi/js';
import Icon from '@mdi/react';
import { useTranslation } from 'react-i18next';

import { DialogModal } from '../../components/ui/DialogModal';
import { List } from '../../components/widget/List';
import type { useSeismicEvents } from './useSeismicEvents';

export const EventListModal = ({
    currentLocale,
    events
}: {
    currentLocale: string;
    events: ReturnType<typeof useSeismicEvents>;
}) => {
    const { t } = useTranslation();
    const {
        selectedSeismicEventSource,
        isSelectModalOpen,
        setIsSelectModalOpen,
        seismicEventsData,
        handleChooseEvent
    } = events;

    return (
        <DialogModal
            heading={
                <div className="space-y-4 text-gray-800">
                    <div className="flex items-center space-x-2 text-lg font-extrabold">
                        <Icon path={mdiTarget} size={1.2} className="flex-shrink-0" />
                        <h2>{t('views.History.event_list_modal.title')}</h2>
                    </div>
                    <div className="ml-0.5 flex items-center space-x-2 text-sm">
                        <img
                            src={selectedSeismicEventSource?.countryFlag}
                            className="size-5 flex-shrink-0"
                        />
                        <span>
                            {selectedSeismicEventSource?.locales[currentLocale] ??
                                selectedSeismicEventSource?.locales[
                                    selectedSeismicEventSource?.defaultLocale
                                ]}
                        </span>
                    </div>
                </div>
            }
            open={isSelectModalOpen}
            onClose={() => setIsSelectModalOpen(false)}
        >
            <List
                className="overflow-y-scroll"
                onClick={(id) => handleChooseEvent(id)}
                data={Object.values(seismicEventsData).sort((a, b) => b.timestamp - a.timestamp)}
            />
        </DialogModal>
    );
};
