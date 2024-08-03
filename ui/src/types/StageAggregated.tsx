import Stage from "./Stage";
import Story from "./Story";
import Channel from "./Channel";

type StageAggregated = {
    stage: Stage;
    story: Story;
    channel: Channel;
  };
  
  export default StageAggregated;