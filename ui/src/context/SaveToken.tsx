function SaveToken(
  token: string,
  expire: EpochTimeStamp,
  user_id: number,
  username: string
): void {
  sessionStorage.setItem("token", token);
  sessionStorage.setItem("expire", expire.toString());
  sessionStorage.setItem("user_id", user_id.toString());
  sessionStorage.setItem("username", username);
}

export default SaveToken;
