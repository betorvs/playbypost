import { useParams } from "react-router-dom";
import { useContext, useEffect, useState } from "react";
import { AuthContext } from "../context/AuthContext";
import Layout from "../components/Layout";
import StageDetailHeader from "../components/StageDetailHeader";
import Encounter from "../types/Encounter";
import StageEncounterCards from "../components/Cards/StageEncounter";
import { FetchStage, FetchStageEncountersByIDWithPagination } from "../functions/Stages";
import { useTranslation } from "react-i18next";
import StageAggregated from "../types/StageAggregated";
import ReactPaginate from "react-paginate";

const StageDetail = () => {
  const postsPerPage = 2;
  const { id, story } = useParams();
  const { Logoff } = useContext(AuthContext);
  const [stage, setStage] = useState<StageAggregated>();

  const safeID: string = id ?? "";

  const storySafeID: string = story ?? "";

  const { t } = useTranslation(["home", "main"]);

  const [encounters, setEncounters] = useState<Encounter[]>([]);

  const [cursor, setCursor] = useState<string>("0");
  const [pageCount, setPageCount] = useState(0); // Total number of pages
  const [currentPage, setCurrentPage] = useState(0);

  useEffect(() => {
    FetchStage(safeID, setStage);
    FetchStageEncountersByIDWithPagination(safeID, cursor, postsPerPage, setPageCount, setEncounters, setCursor);
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
        {stage && <StageDetailHeader detail={true} backButtonLink="/stages" stage={stage} />}
        <div className="row mb-2" key="2">
          {encounters.length !== 0 ? (
            encounters.map((encounter, index) => (
              <StageEncounterCards encounter={encounter} key={index} stageID={safeID} storyId={storySafeID} creator_id={stage?.stage.creator_id ?? 0} />
            ))
          ) : (
            <p>{t("stage.no-encounter", {ns: ['main', 'home']})}</p>
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

export default StageDetail;
