# bta-wiki-import

A cli tool to dump various HBS BT "defs" to wikitext.

## Usage

First you need to convert the defs to wikitext or "export" them:
`bta-wiki-import export ./mod-directory ./path-to-wikitext`

Then you can upload them to the wiki or "import" them:
`bta-wiki-import import -u Username@botname -l https://WEBSITE/api.php --passfile ./file-with-password ./path-to-wikitext` 