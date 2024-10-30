type Channel = {
  channel: string;
  active: boolean;
};

type RunningChannels = {
  title: string;
  channel: string;
  kind: string;
}
export default Channel;
export type { RunningChannels };