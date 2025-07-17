import { useEffect, useState } from "react";
import FetchCharacters from "../functions/Characters";
import Players from '../types/Players';
import { useTranslation } from "react-i18next";
import { Link } from "react-router-dom";

const CharactersList = () => {
  const { t } = useTranslation(['home', 'main']);
   
  const [characters, setCharacters] = useState<Players[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const getCharacters = async () => {
      try {
        setIsLoading(true);
        setError(null);
        await FetchCharacters(setCharacters);
      } catch (err) {
        setError(t("common.error-fetching-characters", {ns: ['main', 'home']}));
        console.error(err);
      } finally {
        setIsLoading(false);
      }
    };
    getCharacters();
  }, []);

  if (isLoading) {
    return (
      <div className="container mt-3" key="2">
        <p>{t("common.loading-characters", {ns: ['main', 'home']})}</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="container mt-3" key="2">
        <p className="text-danger">{error}</p>
      </div>
    );
  }

  return (
    <div className="container mt-3" key="2">
      {characters.length === 0 ? (
        <p>{t("common.no-characters-found", {ns: ['main', 'home']})}</p>
      ) : (
        <table className="table table-striped">
          <thead>
            <tr>
              <th>ID</th>
              <th>{t("common.name", {ns: ['main', 'home']})}</th>
              <th>{t("common.status", {ns: ['main', 'home']})}</th>
              <th>{t("common.actions", {ns: ['main', 'home']})}</th>
            </tr>
          </thead>
          <tbody>
            {characters.map((character) => (
              <tr key={character.id}>
                <td>{character.id}</td>
                <td>{character.name}</td>
                <td>{character.destroyed ? t("common.inactive", {ns: ['main', 'home']}) : t("common.active", {ns: ['main', 'home']})}</td>
                <td>
                  <Link to={`/characters/${character.id}/edit`} className="btn btn-primary btn-sm">
                    {t("common.edit", {ns: ['main', 'home']})}
                  </Link>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
};

export default CharactersList;