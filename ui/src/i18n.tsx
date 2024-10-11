import i18next from "i18next";
import { initReactI18next } from "react-i18next";

//Import all translation files
import translationEnglish from "./translation/english/translation.json";
import translationBRPortuguese from "./translation/br-portuguese/translation.json";

const resources = {
    en: {
        home: translationEnglish,
    },
    pt: {
        home: translationBRPortuguese,
    },
}

i18next
.use(initReactI18next)
.init({
  resources,
  lng:"en", //default language
});

export default i18next;