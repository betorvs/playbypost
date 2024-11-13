import { useParams } from "react-router-dom";
import EncounterCards from "../components/Cards/Encounter";
import { useContext, useEffect, useState } from "react";
import Encounter from "../types/Encounter";
import { AuthContext } from "../context/AuthContext";
import Layout from "../components/Layout";
import StoryDetailHeader from "../components/StoryDetailHeader";
import {FetchEncountersWithPagination} from "../functions/Encounters";
import { useTranslation } from "react-i18next";
import ReactPaginate from "react-paginate";

const StoryDetail = () => {
  const postsPerPage = 4;
  const { id } = useParams();
  const { Logoff } = useContext(AuthContext);
  const { t } = useTranslation(['home', 'main']);

  const safeID: string = id ?? "";

  const [encounters, setEncounters] = useState<Encounter[]>([]);

  const [cursor, setCursor] = useState<string>("0");
  const [pageCount, setPageCount] = useState(0); // Total number of pages
  const [currentPage, setCurrentPage] = useState(0);

  useEffect(() => {
    FetchEncountersWithPagination(safeID, cursor, postsPerPage, setPageCount, setEncounters, setCursor);
  }, [currentPage]);

  const handlePageClick = (data: { selected: number }) => {
    console.log("Selected: " + data.selected);
    setCurrentPage(data.selected);
    const offset = Math.ceil(postsPerPage * data.selected);
    setCursor(offset.toString());
  }

  return (
    <>
      <div className="container mt-3" key="1">
        <Layout Logoff={Logoff} />
        {<StoryDetailHeader detail={true} id={safeID} />}

        <div className="row mb-2" key="2">
          {encounters.length !== 0 ? (
            encounters.map((encounter, index) => (
              <EncounterCards encounter={encounter} key={index} disable_footer={false} />
            ))
          ) : (
            <p>{t("story.error", {ns: ['main', 'home']})}</p>
          )}
          <ReactPaginate
            previousLabel={t("pagination.previous", {ns: ['main', 'home']})}
            nextLabel={t("pagination.next", {ns: ['main', 'home']})}
            breakLabel={"..."}
            breakClassName={"break-me"}
            pageCount={pageCount}
            marginPagesDisplayed={1}
            pageRangeDisplayed={2}
            onPageChange={handlePageClick}
            containerClassName={"pagination"}
            activeClassName={"active"}
          />
        </div>
      </div>
    </>
  );
};

export default StoryDetail;
