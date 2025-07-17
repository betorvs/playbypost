import { useContext } from 'react';
import { useTranslation } from 'react-i18next';
// import FetchCharacters from '../functions/Characters';
// import CharacterCard from '../types/CharacterCard';
import Layout from '../components/Layout';
import { AuthContext } from '../context/AuthContext';
import CharactersList from '../components/CharactersList';

const CharacterListPage = () => {
  const { Logoff } = useContext(AuthContext);
  const { t } = useTranslation(['home', 'main']);
  // const [characters, setCharacters] = useState<CharacterCard[]>([]);

    // useEffect(() => {
    //   FetchCharacters(setCharacters);
    // }, []);

  return (
    <>
    <div className="container mt-3" key="1">
      <Layout Logoff={Logoff} />
      <h1>{t("common.character-list", {ns: ['main', 'home']})}</h1>
      <hr />
    </div>
    {<CharactersList />}
    </>
  );
};

export default CharacterListPage;
