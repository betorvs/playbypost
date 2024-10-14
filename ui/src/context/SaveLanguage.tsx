
function SaveLanguage(language: string) {
    sessionStorage.setItem("language", language);
}

function GetLanguage(): string {
    let token = sessionStorage.getItem("language") || "en";
    return token;
  }

export default SaveLanguage;
export { GetLanguage}; 