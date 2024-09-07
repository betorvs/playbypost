package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/betorvs/playbypost/core/sys/web/types"
	// Add this import statement
)

func (a MainApi) GetStage(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	obj, err := a.db.GetStage(a.ctx)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "stage database issue")
		return
	}
	a.s.JSON(w, obj)
}

func (a MainApi) CreateStage(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	obj := types.Stage{}
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		a.logger.Error("json error ", "error", err.Error())
		a.s.ErrJSON(w, http.StatusBadRequest, "invalid json")
		return
	}
	if obj.StoryID == 0 || obj.UserID == "" {
		a.s.ErrJSON(w, http.StatusBadRequest, "text and story_id cannot be empty")
		return
	}
	res, err := a.db.CreateStageTx(a.ctx, obj.Text, obj.UserID, obj.StoryID)
	if err != nil {
		a.logger.Error("error ", "story_id", obj.StoryID)
		a.s.ErrJSON(w, http.StatusBadGateway, "error creating stage on database")
		return
	}
	msg := fmt.Sprintf("stage id %v", res)
	a.s.JSON(w, types.Msg{Msg: msg})
}

func (a MainApi) GetStageById(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	// headerUsername := r.Header.Get(types.HeaderUsername)
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "id should be a integer")
		return
	}
	obj, err := a.db.GetStageByStageID(a.ctx, id)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "stage issue")
		return
	}
	// if a.Sessions.Current[headerUsername].UserID != obj.Stage.StorytellerID {
	// 	a.s.JSON(w, "{}")
	// 	return
	// }
	a.s.JSON(w, obj)
}

func (a MainApi) GetStageByStoryId(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	// headerUsername := r.Header.Get(types.HeaderUsername)
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "id should be a integer")
		return
	}
	obj, err := a.db.GetStageByStoryID(a.ctx, id)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "stage issue")
		return
	}
	// if a.Sessions.Current[headerUsername].UserID != obj.Stage.StorytellerID {
	// 	a.s.JSON(w, "{}")
	// 	return
	// }
	a.s.JSON(w, obj)
}

func (a MainApi) AddEncounterToStage(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	obj := types.EncounterAssociation{}
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		a.logger.Error("json error ", "error", err.Error())
		a.s.ErrJSON(w, http.StatusBadRequest, "invalid json")
		return
	}
	stage, err := a.db.GetStageByStageID(a.ctx, obj.StageID)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "stage issue")
		return
	}
	if stage.Stage.ID < 0 {
		a.s.ErrJSON(w, http.StatusBadRequest, "stage not found")
		return
	}
	if stage.Stage.StoryID != obj.StoryID {
		a.s.ErrJSON(w, http.StatusBadRequest, "story id does not match with stage story id")
		return
	}

	res, err := a.db.AddEncounterToStage(a.ctx, obj.Text, obj.StageID, stage.Stage.StorytellerID, obj.EncounterID)
	if err != nil {
		a.logger.Error("error ", "stage_id", obj.StageID)
		a.s.ErrJSON(w, http.StatusBadGateway, "error adding encounter to stage on database")
		return
	}
	msg := fmt.Sprintf("stage_encounter id %v", res)
	a.s.JSON(w, types.Msg{Msg: msg})
}

func (a MainApi) AddChannelToStage(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	obj := types.Channel{}
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		a.logger.Error("json error ", "error", err.Error())
		a.s.ErrJSON(w, http.StatusBadRequest, "invalid json")
		return
	}
	stage, err := a.db.GetStageByStageID(a.ctx, obj.StageID)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "stage issue")
		return
	}
	if stage.Stage.ID < 0 {
		a.s.ErrJSON(w, http.StatusBadRequest, "stage not found")
		return
	}

	res, err := a.db.AddChannelToStage(a.ctx, obj.Channel, obj.StageID)
	if err != nil {
		a.logger.Error("error ", "stage_id", obj.StageID)
		a.s.ErrJSON(w, http.StatusBadGateway, "error adding channel to stage on database")
		return
	}
	msg := fmt.Sprintf("stage_channel id %v", res)
	event := fmt.Sprintf("Title: %s\nAnnounce: %s\n", stage.Story.Title, stage.Story.Announcement)
	a.logger.Info("event content", "event", event)
	resp, err := a.client.PostEvent(obj.Channel, "ALL", event, types.EventAnnounce)
	if err != nil {
		a.logger.Error("cannot post to slack", "error", err.Error())
		// msg = fmt.Sprintf("stage_channel id %v but message not posted to slack", res)
		a.s.ErrJSON(w, http.StatusBadGateway, "error sending message to slack")
		return
	}
	a.logger.Info("add channel okay", "body", string(resp))
	a.s.JSON(w, types.Msg{Msg: msg})
}

func (a MainApi) GetStageEncounterByStageID(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "id should be a integer")
		return
	}
	obj, err := a.db.GetStageEncountersByStageID(a.ctx, id)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "stage_encounters database issue")
		return
	}
	a.logger.Info("encounter list", "obj", obj)
	a.s.JSON(w, obj)
}

func (a MainApi) UpdateEncounterPhaseById(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "id should be a integer")
		return
	}
	phaseString := r.PathValue("phase")
	phase, err := strconv.Atoi(phaseString)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "phase should be a integer")
		return
	}
	err = a.db.UpdatePhase(a.ctx, id, phase)
	if err != nil {
		a.logger.Error("error from db", "error", err.Error())
		errMsg := "encounters issue"
		if strings.HasPrefix(err.Error(), "stage_id") {
			errMsg = err.Error()
		}
		a.s.ErrJSON(w, http.StatusBadRequest, errMsg)
		return
	}
	status := types.PhaseAtoi(phase)
	a.logger.Info("change phase worked", "phase", status)
	a.s.JSON(w, types.Msg{Msg: fmt.Sprintf("change to phase: %s", status)})
}

func (a MainApi) GetStageEncounterByEncounterID(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "id should be a integer")
		return
	}
	obj, err := a.db.GetStageEncounterByEncounterID(a.ctx, id)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "stage_encounters database issue")
		return
	}
	a.logger.Info("encounter", "obj", obj)
	a.s.JSON(w, obj)
}

func (a MainApi) AddParticipants(w http.ResponseWriter, r *http.Request) {
	obj := types.Participants{}
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "json decode error")
		return
	}
	if len(obj.Identifies) == 0 || obj.EncounterID == 0 {
		a.s.ErrJSON(w, http.StatusBadRequest, "players id list and encounter id cannot be empty")
		return
	}
	err = a.db.AddParticipants(a.ctx, obj.EncounterID, obj.NPC, obj.Identifies)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "error adding participants to encounter on database")
		return
	}
	msg := fmt.Sprintf("encounter id %v participants updated", obj.EncounterID)
	a.s.JSON(w, types.Msg{Msg: msg})
}

// stage_next_encounter
func (a MainApi) AddNextEncounter(w http.ResponseWriter, r *http.Request) {
	obj := types.NextEncounter{}
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "json decode error")
		return
	}
	if obj.NextEncounterID == 0 || obj.EncounterID == 0 || obj.StageID == 0 {
		a.s.ErrJSON(w, http.StatusBadRequest, "next encounter id, encounter id and stage id cannot be empty")
		return
	}
	if obj.Objective.Kind == "" {
		obj.Objective.Kind = types.ObjectiveDefault
		obj.Objective.Values = []int{0}
	}
	err = a.db.AddNextEncounter(a.ctx, obj)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "error adding next encounter to encounter on database")
		return
	}
	msg := fmt.Sprintf("encounter id %v next encounter updated", obj.EncounterID)
	a.s.JSON(w, types.Msg{Msg: msg})
}

// stage_running_tasks
func (a MainApi) AddRunningTask(w http.ResponseWriter, r *http.Request) {
	obj := types.RunningTask{}
	err := json.NewDecoder(r.Body).Decode(&obj)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "json decode error")
		return
	}
	if obj.StageID == 0 || obj.TaskID == 0 || obj.StorytellerID == 0 || obj.EncounterID == 0 {
		a.s.ErrJSON(w, http.StatusBadRequest, "stage id, task id, storyteller id and encounter id cannot be empty")
		return
	}
	err = a.db.AddRunningTask(a.ctx, obj.Text, obj.StageID, obj.TaskID, obj.StorytellerID, obj.EncounterID)
	if err != nil {
		a.logger.Error("error adding running task", "error", err.Error())
		a.s.ErrJSON(w, http.StatusBadRequest, "error adding running task to encounter on database")
		return
	}
	msg := fmt.Sprintf("task id %v running task updated", obj.TaskID)
	a.s.JSON(w, types.Msg{Msg: msg})
}

func (a MainApi) GetStageEncounterActivitiesByEncounterID(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "id should be a integer")
		return
	}
	activities, err := a.db.GetStageEncounterActivitiesByEncounterID(a.ctx, id)
	if err != nil {
		a.logger.Error("error getting encounter activities", "error", err.Error())
		a.s.ErrJSON(w, http.StatusBadRequest, "error getting encounter activities to encounter on database")
		return
	}
	a.s.JSON(w, activities)
}

func (a MainApi) GetStageEncounterActivities(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	activities, err := a.db.GetStageEncounterActivities(a.ctx)
	if err != nil {
		a.logger.Error("error getting encounter activities", "error", err.Error())
		a.s.ErrJSON(w, http.StatusBadRequest, "error getting encounter activities to encounter on database")
		return
	}
	a.s.JSON(w, activities)
}

func (a MainApi) CloseStage(w http.ResponseWriter, r *http.Request) {
	if a.Session.CheckAuth(r) {
		a.s.ErrJSON(w, http.StatusForbidden, "required authentication headers")
		return
	}
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		a.s.ErrJSON(w, http.StatusBadRequest, "id should be a integer")
		return
	}
	err = a.db.CloseStage(a.ctx, id)
	if err != nil {
		a.logger.Error("error closing stage", "error", err.Error())
		a.s.ErrJSON(w, http.StatusBadRequest, "error closing stage on database")
		return
	}
	a.s.JSON(w, types.Msg{Msg: fmt.Sprintf("stage id %v closed", id)})
}
