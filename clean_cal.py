import os
import pickle
import datetime
from googleapiclient.discovery import build
from google_auth_oauthlib.flow import InstalledAppFlow
from google.auth.transport.requests import Request

SCOPES = ['https://www.googleapis.com/auth/calendar']

calIds = []

def authenticate_google_calendar():
    """Authenticate with Google Calendar API and return a service object."""
    creds = None
    if os.path.exists('env/token.pickle'):
        with open('env/token.pickle', 'rb') as token:
            creds = pickle.load(token)
    if not creds or not creds.valid:
        if creds and creds.expired and creds.refresh_token:
            creds.refresh(Request())
        else:
            flow = InstalledAppFlow.from_client_secrets_file(
                'env/credentials.json', SCOPES)
            creds = flow.run_local_server(port=0)
        with open('env/token.pickle', 'wb') as token:
            pickle.dump(creds, token)
    
    return build('calendar', 'v3', credentials=creds)

def delete_all_events(calendar_id, service):
    events_result = service.events().list(calendarId=calendar_id).execute()
    events = events_result.get('items', [])
    
    if not events:
        print('No upcoming events found.')
    else:
        for event in events:
            try:
                service.events().delete(calendarId=calendar_id, eventId=event['id']).execute()
                print(f"Deleted event: {event['summary']}")
            except Exception as e:
                print(f"An error occurred: {e}")

if __name__ == '__main__':
    service = authenticate_google_calendar()
    for calId in calIds:
        delete_all_events(calId, service)
