import os
import pickle
import datetime
from datetime import datetime, timedelta
from googleapiclient.discovery import build
from google_auth_oauthlib.flow import InstalledAppFlow
from google.auth.transport.requests import Request
from googleapiclient.errors import HttpError

SCOPES = ['https://www.googleapis.com/auth/calendar']

calIds = [
    "aa7dfbd43b5f0ffa3a63a57faae946d8bd317f3c45bfe221919ce5f7c90852fc@group.calendar.google.com",
"567bfc5ca2de6bc54583aa4744cf049a65c451d2bd897d9dc003aa9dfa287a87@group.calendar.google.com",
"59bfd0df6be1efe40b5f10c45d663cd1d14b924ebf3b709cfcd69fb09897138f@group.calendar.google.com",
"f1e41b8ff60c54795946e150920b2971cdf929a015fee172824fb0e961accd42@group.calendar.google.com",
"efee209e8160bf8bb891ab36abc86d7e3a8ea82dd2fd6a7d5cf81ac1a61fb42d@group.calendar.google.com",
"e9c080d310635936318c85aa7c308baabbfb251eabd05342e52f47e20eae4816@group.calendar.google.com",
"73982ef0ab637a3fc9d5645be9018b330669937add99b6825c7a09725bbe6df8@group.calendar.google.com",
"6d8844a5a5e8b4bb4fcc83c7e0571640e11ff7c2e2f1989f79bee62fda5f57cb@group.calendar.google.com",
"6bf6d9cce226a2072bbbf0c16a0c0f96af99c424d167f353d9ae0b57eb6da23f@group.calendar.google.com",
"9ca57980980ce6b06a929b6eedbc624b8a1d3227733fb86946d65242e5752ae4@group.calendar.google.com",
]

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

def delete_events_after_date(calendar_id, service, after_date):
    service = authenticate_google_calendar()
    
    after_date_rfc3339 = after_date.isoformat() + 'Z'
    
    events_result = service.events().list(
        calendarId=calendar_id,
        timeMin=after_date_rfc3339,
        singleEvents=True,
        orderBy='startTime'
    ).execute()
    
    events = events_result.get('items', [])
    
    if not events:
        print('No events found.')
        return

    
    for event in events:
        try:
            event['summary']
        except Exception:
            print(f'No summary found for event ID: {event["id"]}')
            service.events().delete(calendarId=calendar_id, eventId=event['id']).execute()
            continue
        if not event['summary'].startswith("MUSbooking"):
            print(f'Deleting event: {event["summary"]} (ID: {event["id"]})')
            service.events().delete(calendarId=calendar_id, eventId=event['id']).execute()


if __name__ == '__main__':
    service = authenticate_google_calendar()
    for calId in calIds:
        delete_events_after_date(calId, service, datetime(2024, 9, 28))
