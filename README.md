# Nextcloud-TrelloSync
 Service to syncronize from Nextcloud to Trello based on a simple config.

# Features
Since this was made on a whim within 3-4 hours one night, it only has a very basic feature set to *replicate* the cards from Nextcloud Deck to Trello.

There is also a *explorer* included that uses the credentials in your config file to get the different IDs for Nextcloud Deck and Trello respectively.

The explorer program defaults to showing both Nextcloud Deck and Trello info, or it can take one argument which is either nextcloud, deck or trello.

There are still several API calls for both of them created (as seen in [pkg/api](pkg/api)):
- [x] Get boards
- [x] Get stacks/lists
- [x] Get single stack/list
- [x] Create new card
- [x] Delete existing card
- [x] Update card (Trello only atm.)

## To-Do
- [ ] Actually move cards instead of deleting and recreating it - do this based on the name
- [ ] Archiving cards
- [ ] Create Trello to Nextcloud Deck sync
- [ ] More API calls?

# Config
The config is *rather* simple, here is an example:
```
debug: true
log: debug.log
interval: 10
nextcloud:
  api: https://cloud.example.com/index.php/apps/deck/api/v1.0
  username: Mikkel
  password: 1234
trello:
  key: 1234
  token: 4321

sync:
- name: "To-Do"
  type: "trello"
  nextcloud:
    board: 17
    stack: 72
  trello:
    board: 5fe970e7eb778e3c11308bda
    list: 5fe999817406ae352549964c
```
- `debug` - show debug in logs or not.
- `log` - path to write log-file to.
- `interval` - in minutes for how often it should try synchronizing.
- `nextcloud` - Nextcloud credenticals and API endpoint, note that the password is likely a device password, [read more](https://docs.nextcloud.com/server/stable/user_manual/en/session_management.html#managing-devices).
- `trello` - your Trello API key and account token respectively, [read more](https://developer.atlassian.com/cloud/trello/guides/rest-api/api-introduction/).
- `sync` - the *synchronization groups* which is where you configure the board and stack/list to actually synchronize.
  - `name` - name for this group, is only used for debugging and logs.
  - `type` - which service to syncronize **to** - *trello* or *nextcloud* (not implemented yet).
  - The rest are just the IDs for the board and stack/list to synchronize.