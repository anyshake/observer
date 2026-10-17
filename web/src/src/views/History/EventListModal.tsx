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
                <div className="app-dialog-intro">
                    <div className="app-dialog-heading">
                        <Icon path={mdiTarget} size={1.2} className="app-icon" />
                        <h2>{t('views.History.event_list_modal.title')}</h2>
                    </div>
                    <div className="app-event-subtitle">
                        <img
                            src={selectedSeismicEventSource?.countryFlag}
                            className="app-event-flag"
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
                className="app-scroll-y"
                onClick={(id) => handleChooseEvent(id)}
                data={Object.values(seismicEventsData).sort((a, b) => b.timestamp - a.timestamp)}
            />
        </DialogModal>
    );
};
