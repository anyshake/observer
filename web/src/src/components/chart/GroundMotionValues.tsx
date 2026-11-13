import { useTranslation } from 'react-i18next';

import { formatGroundMotion, IGroundMotion } from './groundMotion';

interface IGroundMotionValues {
    readonly motion?: IGroundMotion | null;
    readonly className?: string;
}

export const GroundMotionValues = ({ motion, className }: IGroundMotionValues) => {
    const { t } = useTranslation();
    const items = [
        { key: 'pga', value: motion?.pga, unit: t('components.GroundMotion.pga_unit') },
        { key: 'pgv', value: motion?.pgv, unit: t('components.GroundMotion.pgv_unit') },
        { key: 'pgd', value: motion?.pgd, unit: t('components.GroundMotion.pgd_unit') }
    ];
    return (
        <dl
            className={`app-ground-motion ${className ?? ''}`}
            title={t('components.GroundMotion.description')}
        >
            {items.map(({ key, value, unit }) => (
                <div key={key} className="app-ground-motion-item">
                    <dt className="app-ground-motion-label">
                        {t(`components.GroundMotion.${key}`)}
                    </dt>
                    <dd className="app-ground-motion-value">
                        {formatGroundMotion(value)}
                        <span className="app-ground-motion-unit">{unit}</span>
                    </dd>
                </div>
            ))}
        </dl>
    );
};
