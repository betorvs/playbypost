// type SessionToken = {
//   access_token: string;
//   expire_on: EpochTimeStamp;
//   user_id: number;
//   status: string;
//   msg: string;
// };

export interface SessionToken {
  access_token: string;
  expire_on: EpochTimeStamp;
  user_id: number;
  status: string;
  msg: string;
}

export interface Session {
  Username: string;
  Token: string;
  UserID: number;
  Expiry: string;
  ClientType: string;
  ClientInfo: string;
  IPAddress: string;
  UserAgent: string;
  CreatedAt: string;
  UpdatedAt: string;
  LastActivity: string;
}

export interface SessionEvent {
  ID: number;
  SessionID: number;
  EventType: string;
  Timestamp: string;
  EventData: any; // Use 'any' for now, can be more specific later
}
