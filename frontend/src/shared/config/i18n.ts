import i18n from 'i18next'
import { initReactI18next } from 'react-i18next'

const resources = {
    ru: {
        translation: {
            nav: {
                home: 'Главная',
                me: 'Мои события',
                organizer: 'Консоль организатора',
                logout: 'Выйти',
            },
            common: {
                loading: 'Загрузка...',
                notFound: 'Страница не найдена',
            },
        },
    },
    en: {
        translation: {
            nav: {
                home: 'Home',
                me: 'My events',
                organizer: 'Organizer console',
                logout: 'Logout',
            },
            common: {
                loading: 'Loading...',
                notFound: 'Page not found',
            },
        },
    },
}

i18n.use(initReactI18next).init({
    resources,
    lng: 'ru',
    fallbackLng: 'en',
    interpolation: {
        escapeValue: false,
    },
})

export default i18n
