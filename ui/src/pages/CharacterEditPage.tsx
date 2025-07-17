import { useContext, useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import Layout from '../components/Layout';
import { AuthContext } from '../context/AuthContext';
import FetchCharacters from '../functions/Characters';
import Players from '../types/Players';
import UpdateCharacter from '../functions/UpdateCharacter';

const CharacterEditPage = () => {
  const { Logoff } = useContext(AuthContext);
  const { t } = useTranslation(['home', 'main']);

  const [characters, setCharacters] = useState<Players[]>([]);
  const [selectedCharacter, setSelectedCharacter] = useState<Players | null>(null);
  const [characterName, setCharacterName] = useState('');
  const [characterRPG, setCharacterRPG] = useState('');
  const [nameError, setNameError] = useState<string | null>(null);
  const [rpgError, setRpgError] = useState<string | null>(null);

  useEffect(() => {
    FetchCharacters(setCharacters);
  }, []);

  useEffect(() => {
    if (selectedCharacter) {
      setCharacterName(selectedCharacter.name);
      setCharacterRPG(selectedCharacter.rpg);
    } else {
      setCharacterName('');
      setCharacterRPG('');
    }
  }, [selectedCharacter]);

  const handleCharacterChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const characterId = parseInt(e.target.value);
    const character = characters.find(char => char.id === characterId);
    setSelectedCharacter(character || null);
  };



  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    // Reset errors
    setNameError(null);
    setRpgError(null);

    let isValid = true;

    if (!selectedCharacter) {
      console.error("No character selected.");
      isValid = false;
    }

    if (characterName.trim() === '') {
      setNameError(t("common.character-name-required", {ns: ['main', 'home']}));
      isValid = false;
    }

    if (characterRPG.trim() === '') {
      setRpgError(t("common.character-rpg-required", {ns: ['main', 'home']}));
      isValid = false;
    }

    if (!isValid) {
      return;
    }

    const updatedData = {
      name: characterName,
      rpg: characterRPG,
      // Add other fields as they become editable
    };

    const success = await UpdateCharacter(selectedCharacter!.id, updatedData);
    if (success) {
      alert(t("common.character-updated-successfully", {ns: ['main', 'home']}));
      // Optionally, refresh the character list or navigate away
    } else {
      alert(t("common.failed-to-update-character", {ns: ['main', 'home']}));
    }
  };

  return (
    <>
      <div className="container mt-3">
        <Layout Logoff={Logoff} />
        <h1>{t("common.edit-character", {ns: ['main', 'home']})}</h1>
        <hr />
        <form onSubmit={handleSubmit}>
          <div className="mb-3">
            <label htmlFor="characterSelect" className="form-label">
              {t("common.select-character", {ns: ['main', 'home']})}
            </label>
            <select
              className="form-select"
              id="characterSelect"
              onChange={handleCharacterChange}
              value={selectedCharacter?.id || ''}
            >
              <option value="">{t("common.select-a-character", {ns: ['main', 'home']})}</option>
              {characters.map((character) => (
                <option key={character.id} value={character.id}>
                  {character.name}
                </option>
              ))}
            </select>
          </div>

          {selectedCharacter && (
            <>
              <div className="mb-3">
                <label htmlFor="characterName" className="form-label">
                  {t("common.character-name", {ns: ['main', 'home']})}
                </label>
                <input
                  type="text"
                  className={`form-control ${nameError ? 'is-invalid' : ''}`}
                  id="characterName"
                  value={characterName}
                  onChange={(e) => setCharacterName(e.target.value)}
                  required
                />
                {nameError && <div className="invalid-feedback">{nameError}</div>}
              </div>
              <div className="mb-3">
                <label htmlFor="characterRPG" className="form-label">
                  {t("common.character-rpg", {ns: ['main', 'home']})}
                </label>
                <input
                  type="text"
                  className={`form-control ${rpgError ? 'is-invalid' : ''}`}
                  id="characterRPG"
                  value={characterRPG}
                  onChange={(e) => setCharacterRPG(e.target.value)}
                  required
                  disabled={true}
                />
                {rpgError && <div className="invalid-feedback">{rpgError}</div>}
              </div>
              <button type="submit" className="btn btn-primary">
                {t("common.save-changes", {ns: ['main', 'home']})}
              </button>
            </>
          )}
        </form>
      </div>
    </>
  );
};

export default CharacterEditPage;
