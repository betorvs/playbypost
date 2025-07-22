import { useEffect, useState } from 'react';
import { Session, SessionEvent } from '../types/Session.tsx';
import { fetchAllSessions, fetchSessionEvents } from '../functions/SessionApi';
import { formatDateTime } from '../functions/time.ts';
import { useTranslation } from 'react-i18next';

const SessionMonitoring = () => {
    const [sessions, setSessions] = useState<Session[]>([]);
    const [sessionEvents, setSessionEvents] = useState<SessionEvent[]>([]);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);
    const { t } = useTranslation(['main', 'home']);
  
    useEffect(() => {
      const loadData = async () => {
        try {
          setLoading(true);
          const fetchedSessions = await fetchAllSessions();
          // console.log('Fetched Sessions:', fetchedSessions);
          setSessions(fetchedSessions);
  
          const fetchedEvents = await fetchSessionEvents();
          setSessionEvents(fetchedEvents);
        } catch (err) {
          setError('Failed to load session data.');
          console.error(err);
        } finally {
          setLoading(false);
        }
      };
  
      loadData();
    }, []);
  
    if (loading) {
      return <div>Loading session data...</div>;
    }
  
    if (error) {
      return <div>Error: {error}</div>;
    }
  return <div className="container mt-3" key="2">

  <section>
    <h2>{t("common.active-sessions", {ns: ['main', 'home']})}</h2>
    {sessions.length === 0 ? (
      <p>{t("common.no-active-sessions", {ns: ['main', 'home']})}</p>
    ) : (
      <table className="table table-striped">
        <thead>
          <tr>
            <th>{t("common.username", {ns: ['main', 'home']})}</th>
            <th>{t("common.client-type", {ns: ['main', 'home']})}</th>
            <th>{t("common.ip-address", {ns: ['main', 'home']})}</th>
            <th>{t("common.last-activity", {ns: ['main', 'home']})}</th>
          </tr>
        </thead>
        <tbody>
          {sessions.map((session) => (
            <tr key={session.Token}>
              <td>{session.Username}</td>
              <td>{session.ClientType}</td>
              <td>{session.IPAddress}</td>
              <td>{formatDateTime(session.LastActivity)}</td>
            </tr>
          ))}
        </tbody>
      </table>
    )}
  </section>

  <section>
    <h2>{t("common.session-events", {ns: ['main', 'home']})}</h2>
    {sessionEvents.length === 0 ? (
      <p>{t("common.no-session-events-found", {ns: ['main', 'home']})}</p>
    ) : (
      <table>
        <thead>
          <tr>
            <th>{t("common.event-type", {ns: ['main', 'home']})}</th>
            <th>{t("common.timestamp", {ns: ['main', 'home']})}</th>
            <th>{t("common.session-id", {ns: ['main', 'home']})}</th>
            <th>{t("common.event-data", {ns: ['main', 'home']})}</th>
          </tr>
        </thead>
        <tbody>
          {sessionEvents.map((event) => (
            <tr key={event.ID}>
              <td>{event.EventType}</td>
              <td>{formatDateTime(event.Timestamp)}</td>
              <td>{event.SessionID}</td>
              <td>{JSON.stringify(event.EventData)}</td>
            </tr>
          ))}
        </tbody>
      </table>
    )}
  </section>
</div>;
};

export default SessionMonitoring;