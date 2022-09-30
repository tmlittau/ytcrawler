# ytcrawler
Web Crawler to pull different statistics from YouTube.

## Current Implementation *prototype will be removed in future builds*
Currently the jupyter notebook prototype takes a given channel id and pulls the statistics of the channels the given channel is subscribed to.

E.g. Given your own channel id, the prototype pulls the statistics of all the channels you are subscribed to.

Currently only total views, subscriber count and number of videos is pulled. This is visualized in a blob diagram.


### Web Backend 

#### Dependencies
youtube api

  `go get -u google.golang.org/api/youtube/v3`
  
oauth2

  `go get -u golang.org/x/oauth2/...`
 

#### Usage Manual
The Backend written in Go currently works in the same way as the jupyter prototype.

To use the Script, OAuth 2.0 is required. The following instruction on OAuth is taken from https://developers.google.com/youtube/v3/quickstart/go

  1. Use this wizard to create or select a project in the Google Developers Console and automatically turn on the API. Click Continue, then Go to credentials.
  2. On the Create credentials page, click the Cancel button.
  3. At the top of the page, select the OAuth consent screen tab. Select an Email address, enter a Product name if not already set, and click the Save button.
  4. Select the Credentials tab, click the Create credentials button and select OAuth client ID.
  5. Select the application type Other, enter the name "YouTube Data API Quickstart", and click the Create button.
  6. Click OK to dismiss the resulting dialog.
  7. Click the file_download (Download JSON) button to the right of the client ID.
  8. Move the downloaded file to your working directory and rename it client_secret.json.
