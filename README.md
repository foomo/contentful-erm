# gocontentful

A Contentful API code generator for Go. Initial features:

- Creation of Value Objects from Contentful content model
- CDA, CPA, CMA support for CRUD operations
- Automatic management/resolution of references and parents
- Many utility functions, like converting to/from RichText to HTML and easy access to Assets

Rationale
-----------------------------------------

While an unofficial/experimental Go client by Contentful Labs has been available for a long time (see credits at the end of this document), working with the response of a JSON REST API is far from optimal. 

Loading and converting data into (and from) Go value objects can quickly become a nightmare and the amount of code to access each field of each content type of a complex space can literally explode and if you're not careful how and when you interact with an API over the Internet, performance can be impacted dramatically. Not to mention what happens to your code if you need to make a significant change the space data model. 

Also, the response of the REST API is very nicely modeled for a dynamic language where _any type_ can be returned in a slice of objects. The strict typing of Go makes this tricky to manage and you end up writing (and re-writing) a lot of similar code. 

The scenario above suggested the need for a code generator that can scan a Contentful space and create a complete, native Go API to interact with the remote service and provide code that can be regenerated in some seconds when the data model change.

How much code is that? As an example of a real-world scenario, a space content model with 11 content types with ranging from 3 to over 40 fields each generated 14,500 lines of Go code. Do you need all those lines? Yes, you do, otherwise it means you never need to read/write some of the content types and fields you defined in the model. In other words, your model needs some spring cleaning.   

Quickstart
-----------------------------------------

### Installation

Prerequisite: you need Go 1.16+. Upgrade if you still haven't, then run:

> go get github.com/foomo/gocontentful

Test the installation (make sure $GOPATH/bin is in your $PATH):

>gocontentful

<pre><code>Contentful API Generator starting...

ERROR: Please specify the Contentful space ID and access Key

SYNOPSIS
     gocontentful -spaceid SpaceID -cmakey CMAKey [-contenttypes firsttype,secondtype...lasttype] path/to/target/package

Usage of gocontentful:
  -cmakey string
    	Contentful CMA key
  -contenttypes string
    	[Optional] Content type IDs to parse, comma separated
  -spaceid string
    	Contentful space ID

Note: The last segment of the path/to/target/package will be used as package name
</code></pre>

### Use case

- You want to generate a package named "people" and manipulate entries of content types with ID "person" and "pet".

Run the following:

>gocontentful -spaceid YOUR_SPACE_ID -cmakey YOUR_CMA_API_TOKEN -contenttypes person,pet path/to/your/go/project/folder/people 

The -contenttypes parameter is optional. If not specified, an API for all content types of the space will be generated.

The script will scan the space, download locales and content types and generate the Go API files in the target path:

<pre><code>path/to/your/go/project/folder/people
|-gocontentfulvobase.go
|-gocontentfulvolib_person.go // One file for each content type
|-gocontentfulvolib_pet.go    // One file for each content type
|-gocontentfulvolib.go
|-gocontentfulvo.go
</code></pre>

_Note: Do NOT modify these files! If you change the content model in Contentful you will need to run the generator again and overwrite the files._

### Get a client

The generated files will be in the "people" subdirectory of your project. Your go program can get a Contentful client from them:

`cc, err := people.NewContentfulClient(YOUR_SPACE_ID, people.ClientModeCDA, YOUR_API_KEY, 1000, contentfulLogger, people.LogDebug,false)`

The parameters to pass to NewContentfulClient are:

- *spaceID* (string) 
- *clientMode* (string) supports the constants ClientModeCDA, ClientModeCPA and ClientModeCMA. If you need to operate on multiple APIs (e.g. one for reading and CMA for writing) you need to get two clients
- *clientKey* (string) is your API key (generate one for your API at Contentful)
- *optimisticPageSize* (uint16) is the page size the client will use to download entries from the space for caching. Contentful's default is 100 but you can specify up to 1000: this might get you into an error because Contentful limits the payload response size to 70 KB but the client will handle the error and reduce the page size automatically until it finds a proper value. Hint: using a big page size that always fails is a waste of time and resources because a lot of initial calls will fail, whereas a too small one will not leverage the full download bandwidth. It's a trial-and-error and you need to find the best value for your case. For simple content types you can start with 1000, for very complex ones that include fat fields you might want to get down to 100 or even less.
- *logFn* is a func(fields map[string]interface{}, level int, args ...interface{}) that the client will call whenever it needs to log something. It can be nil if you don't need logging and that will be handled gracefully but it's not recommended. A simple function you can pass that uses the https://github.com/Sirupsen/logrus package might look something like this:
<pre><code>contentfulLogger := func(fields map[string]interface{}, level int, args ...interface{}) {
    switch level {
    case people.LogDebug:
        log.WithFields(fields).Debug(args)
    case people.LogInfo:
        log.WithFields(fields).Info(args)
    case people.LogWarn:
        log.WithFields(fields).Warn(args)
    case people.LogError:
        log.WithFields(fields).Error(args)
    default:
        return
    }
}
</code></pre>
- *logLevel* (int) is the debug level (see function above). Please note that LogDebug is very verbose and even logs when you request a field value but that is not set for the entry.
- *debug* (bool) is the Contentful API client debug switch. If set to *true* it will log on stdout all the CURL calls to Contentful. This is extremely verbose and extremely valuable when something fails in a call to the API because it's the only way to see the REST API response.

### Caching

<pre><code>contentTypes := []string{"person", "pet"}
err = cc.UpdateCache(context, contentTypes, true)
</code></pre>

If your client mode is ClientModeCDA you can ask the client to cache the space (limited to the content types you pass to the UpdateCache function call). The client will download all the entries, convert and store them in the case as native Go value objects. This makes subsequent accesses to the space data an in-memory operation removing all the HTTP overhead you'd normally experience.

The first parameter is the context. If you don't use a context in your application or service just pass _context.TODO()_

The third parameter of UpdateCache toggles asset caching on or off. If you deal with assets you want this to be always on.

The cache update uses 4 workers to speed up the process. This is safe since Contentful allows up to 5 concurrent connections. If you have content types that have a lot of entries, it might make sense to keep them close to each other in the content types slice passed to UpdateCache(), so that they will run in parallel and not one after the other (in case you have more than 4 content types, that is). 

All functions that query the space through ERM are cache-transparent: if a cache is available data will be loaded from there, otherwise it will be sourced from Contentful.

Note that when something changes in the space at Contentful you need to regenerate the cache. This can be done setting up a webhook at Contentful and handling it in your service.

_@TODO: implement support for Contentful's sync API for selective updates of the cache_ 

### Have fun with persons and pets

Load all persons in the space:

<pre><code>persons, err := cc.GetAllPerson() // also consider using the companion
                                  // cc.GetFilteredPerson(query *contentful.Query) 
                                  // to load filtered entries 
</code></pre>
Load a specific person: 

`person, err := cc.GetPersonByID(THE_PERSON_ID)`

ERM provides getters and setters for all fields in the entry. Get that person's name (assuming the entry has a "name" field):

`name := person.Name() // returns Jane`

Get Jane's work title in a different locale:

`name := person.Title(people.SpaceLocaleItalian)`

Note that constants are available for all locales supported by the space. If a space is configured to have a fallback from one locale to the default one, the getter functions will return that in case the value is not set for locale passed to the function.

Contentful supports Rich Text editing and sooner or later you'll want to convert that to HTML:

`htmlText := people.RichTextToHtml(person.Resume(), linkResolver, imageResolver)`

_Note: linkResolver and imageResolver are two functions that resolve URLs for links and attributes for embedded image assets. See API documentation below._

...or the other way around (often used when digesting data from external sources):

`myRichText := HtmlToRichText(htmlSrc)`

ERM supports references out of the box, the internals of how those are managed at the API level are completely transparet. To get a list of Jane's pets (assuming _pets_ is a multiple reference field) just do:

`pets := person.Pets()`

The value returned is a slice *[]*EntryReference* where each element carries the ID (string), the content type (string) and the value object as an *interface{}*. In case of references that can return multiple types this is useful to perform a switch on the content type and assert the type of the returned value object: 

<pre><code>for _, pet := range pets {
  switch pet.ContentType {
  case people.ContentTypeDog: // you have these constants in the generated code
    dog := pet.VO.(*people.Dog)
    // do something with dog
  case people.ContentTypeCat:
    // ...
  }
</code></pre>

Note that if all you have is an ID you can get its content type like this:

`contentType, err := cc.GetContentTypeOfID("XYZ123")`

The inverse direction of references is supported as well. Find all entries that reference a specific dog:

`dog.GetParents() // these, again, will be *[]EntryReference`

If Jane's dog entry has an image field (AKA a reference to an asset), it's easy to access it:

<pre><code>picture := dog.Picture()
if picture != nil {
    theURL := picture.Fields.File.URL
    // ...then do something with this
}
</code></pre>

When you need to change the value of the field, you can use any of the setter functions, e.g.:

`err := dog.SetAge(8)`

but consider the following:

- To save the entry you need to use a client you instantiated with *ClientModeCMA*. Entries retrieved with ClientModeCDA or ClientModeCPA can be saved in memory (for example if you need to enrich the built-in cache) but not persisted to Contentful.
- Make sure you Get and entry right before you manipulate it and upsert it / publish it to Contentful. In case it's been saved by someone else in the meantime, the upsert will fail with a version mismatch error.

To upsert (save) an entry: 

`err := dog.UpsertEntry()`

To publish it (after it's been upserted):

`err := dog.PublishEntry() // change your mind with err := dog.UnpublishEntry()`

Do it in one step:

`err := dog.UpdateEntry() // upserts and publishes`

And finally delete it (but don't delete that lovely puppy please!):

`err := dog.DeleteEntry()`

### Environments support

Not there yet, sorry.

_@TODO: add support for environments_ 

### Unit tests

You'll write your own unit tests using the generated code, but the generator itself will need some good tests. These require a non-trivial set of data and a pass-through HTTP connection to load it from the filesystem instead of the network. It's all in the plan, hang on.

_@TODO: add sample space data and unit tests_ 

Public function set
---------------------

**BASE FUNCTIONS COMMON TO ALL CONTENT TYPES**

>**NewContentfulClient**(spaceID string, clientMode string, clientKey string, optimisticPageSize uint16, logFn func(fields map[string]interface{}, level int, args ...interface{}), logLevel int, debug bool) (*ContentfulClient, error)

Creates a Contentful client, this is the first function you need to call. For usage details please refer to the Quickstart above

---

**FUNCTIONS NAMED AFTER THE CONTENT TYPE**

_For these we're assuming a content type named "Person"._

>**NewCfPerson**() (cfPerson *CfPerson)

Creates a new Person entry. You can manipulate and upsert this later.

>(cc *ContentfulClient) **GetAllPerson**() (voMap map[string]*CfPerson, err error)

Retrieves all Person entries from the client and returnes a map where the key is the ID of the entry and the value is the Go value object for that entry.

>(cc *ContentfulClient) **GetFilteredPerson**(query *contentful.Query) (voMap map[string]*CfPerson, err error) 

Retrieves Person entries matching the specified query.

>(cc *ContentfulClient) **GetPersonByID**(id string) (vo *CfPerson, err error)

Retrieves the Person entry with the specified ID.

---

**REFERENCE CONVERSION AND CONTENT TYPE FUNCTIONS**

>(vo *CfPerson) **ToReference**() (refSys ContentTypeSys) 

Converts a value object into a reference that can be added to a reference field of an entry. Note that functions that retrieve referenced entries return a more flexible and useful _[]*EntryReference_ (see Quickstart above) but to store a reference you need a ContentTypeSys.

>(cc *ContentfulClient) **GetContentTypeOfID**(ID string) (contentType string)

Returns the Contentful content type of an entry ID.

---

**ENTRY FIELD GETTERS**

Field getters are named after the field ID in Contentful and return the proper type. For example, if the Person content type has a Symbol (short text) field named 'Name', this will be the getter:

>(vo *CfPerson) **Name**(locale ...string) (string) 

The locale parameter is optional and if not passed, the function will return the value for the default locale of the space. If the locale is specified and it's not available for the space, an error is returned. If the locale is valid but a value doesn't exist for the field and locale, the function will return the value for the default locale if that's specified as a fallback locale in the space definition in Contentful, otherwise will return an error.

Possible return types are:

- _string_ for fields of types Symbol, Text, Date
- _[]string_ for fields of type List
- _float64_ for fields of type Integer or Number
- _bool_ for fields of type Boolean
- _*ContentTypeSys_ for single reference fields
- _[]*ContentTypeSys_ for multiple reference fields
- _*ContentTypeFieldLocation_ for fields of type Location
- *interface{} for fields of type Object or RichText

If logLevel is set to LogDebug retrieving the value of a field that is not set and so not available in the API response even as a fallback to the default locale will log the event. This can become incredibly verbose, use with care.  

---

**ENTRY FIELD SETTERS  (only available for _ClientModeCMA_)**

Field setters are named after the field ID in Contentful and require to pass in the proper type. See FIELD GETTERS above for a reference. Example:

>(vo *CfPerson) **SetName**(title string, locale ...string) (err error) 

---

**ENTRY WRITE OPERATIONS  (only available for _ClientModeCMA_)**

>(vo *CfPerson) **UpsertEntry**(cc *ContentfulClient) (err error) 

Upserts the entry. This will appear as "Draft" (if it's a new entry) or "Changed" if it's already existing. In the latter case, you will need to retrieve the entry with one of the Manage* functions above to acquire the Sys object that contains the version information. Otherwise the API call will fail with a "Version mismatch" error.

>(vo *CfPerson) **PublishEntry**(cc *ContentfulClient) (err error) 

Publishes the entry. Note that before publishing you will need to retrieve the entry with one of the Manage* functions above to acquire the Sys object that contains the version information. Otherwise the API call will fail with a "Version mismatch" error. This is needed even if you have just upserted the entry with the function above!

>(vo *CfPerson) **UnpublishEntry**(cc *ContentfulClient) (err error) 

Unpublishes the entry. Note that before unpublishing you will need to retrieve the entry with one of the Manage* functions above to acquire the Sys object that contains the version information. Otherwise the API call will fail with a "Version mismatch" error. This is needed even if you have just upserted the entry with the function above!

>(vo *CfPerson) **UpdateEntry**(cc *ContentfulClient) (err error) 

Shortcut function that upserts and publishes the entry. Note that before calling this you will need to retrieve the entry with one of the Manage* functions above to acquire the Sys object that contains the version information. Otherwise the API call will fail with a "Version mismatch" error. Using this shortcut function avoids retrieving the entry twice.

>(vo *CfPerson) **DeleteEntry**(cc *ContentfulClient) (err error) 

Unpublishes and deletes the entry

---
**ASSET FUNCTIONS**

>(cc *ContentfulClient) **GetAllAssets()** (map[string]*contentful.Asset, error)

Retrieve all assets from a space

>(cc *ContentfulClient) **GetAssetByID**(id string) (*contentful.Asset, error)

Retrieve an asset from a space by its ID

>**NewAssetFromURL**(id string, uploadUrl string, imageFileType string, title string, locale ...string) *contentful.Asset

Creates an Asset from an URL of an existing file online (you still need to upsert it later).

>**ToAssetReference**(asset *contentful.Asset) (refSys ContentTypeSys)

Converts the asset to a reference. You need to do this before you add the asset to a reference field of an entry.

>(cc *ContentfulClient) **DeleteAsset**(asset *contentful.Asset) error

Deletes an asset from a space by its ID (only available for _ClientModeCMA_)

---
**UTILITY FUNCTIONS**

>**HtmlToRichText**(htmlSrc string) *RichTextNode

Converts an HTML fragment to a RichTextNode. This is far from complete but useful to migrate data from third-party systems to Contentful. It currently supports headings, paragraphs, hyperlinks, italic and bold tags, horizontal rules, blockquote, ordered and unordered lists, code. Unknown tags are stripped. This function doesn't return any error as it converts the input text into something as good as possible, without  validation.

>**RichTextToHtml**(rt interface{}, linkResolver LinkResolverFunc, imageResolver ImageResolverFunc) (string, error)

Converts an interface representing a Contentful RichText value (usually from a field getter) into HTML. It currently supports all tags except for embedded and inline entries and assets. It takes in two functions to resolve hyperlink URLs (if passed as null this will return the href unchanged) and to derive IMG tag attributes for embedded image assets (you usually want this customized in your application code). The two functions have this signature:

>type LinkResolverFunc func(url string) (resolvedURL string, resolveError error)

>type ImageResolverFunc func(assetID string) (attrs map[string]string, resolveError error)


---

**CONSTANTS**

Each generated content type library file exports a constant with the Contentful ID of the content type itself, for example in _contentful_vo_lib_person.go_:

>const ContentTypePerson = "person"

Constants are available for each locale supported by the space at the time of code generation, e.g.:

>const SpaceLocaleGerman Locale = "de"
>const SpaceLocaleFrench Locale = "fr"
>const defaultLocale Locale = SpaceLocaleGerman

Four levels of logging are supported (even if only partially used at this time):

<pre><code>const (
    LogDebug = 0
    LogInfo  = 1
    LogWarn  = 2
    LogError = 3
)</code></pre>


Dependencies and credits
------------------------

The Go package generated by ERM only relies on one external library:

`https://github.com/foomo/contentful`

That is the raw API client used to interact with Contentful's API. This was originally found at https://github.com/contentful-labs/contentful-go and was forked first by https://github.com/ohookins and then by us (https://github.com/foomo/contentful). You will need the foomo version.
