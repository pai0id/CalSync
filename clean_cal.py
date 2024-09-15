import os
import pickle
import datetime
from googleapiclient.discovery import build
from google_auth_oauthlib.flow import InstalledAppFlow
from google.auth.transport.requests import Request

SCOPES = ['https://www.googleapis.com/auth/calendar']

calIds = ["50c82426236347be525723d287c83a20de8b336689c7bcbdc62c2197249f2d7b@group.calendar.google.com",
"e094039b8c6a95961d567523fdf8586f984c5d445b4131d9c1be62511cbbbe84@group.calendar.google.com",
"5bce8f738c60fc06c5b4eba02c21b38b55f00146c069f8f1447b0dd27e562a8f@group.calendar.google.com",
"9e1ad33080757c217c4648f6d8c28b23f652a2c8419533fe77d93465088b91de@group.calendar.google.com",
"3ddce7e25bee8475a47b21fa7f19e3d6188b8501ee1b0eafe90ce561737e6ef2@group.calendar.google.com",
"3da901d09b590abbbbdb77651ae2e8ff1ace46c7e51d10a32d9509f66c55ac33@group.calendar.google.com",
"5cfd1da735d9c354ac4c45024056289e29c580d6effe60a2139819c24b14bbc9@group.calendar.google.com",
"3e9dca44aec311959f94fadf8bb356c2dec2bdf6c82c44ac8af545357d86a2f8@group.calendar.google.com",
"896e972cd329cc65d4fb0320a1fba1fe3a0e71804e7c64563da54d4558819787@group.calendar.google.com",
"2d62001c42ee5c946a2febf2438f51e3ec3a3530319b23f681ded5797db63c47@group.calendar.google.com"]

def authenticate_google_calendar():
    """Authenticate with Google Calendar API and return a service object."""
    creds = None
    if os.path.exists('token.pickle'):
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

def delete_all_events(calendar_id):
    """Delete all events from a Google Calendar."""
    service = authenticate_google_calendar()
    
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
    for calId in calIds:
        delete_all_events(calId)
