# Luster

Luster is tool for scaping data from Facebook. It currently supports only scraping list of all users who like or follow your personal Facebook page. 

## How to install

Download archive for your platform from [releases page](https://github.com/zladovan/luster/releases/latest) and unpack it to some directory on your file system.

### Install from source

Alternatively you can install `luster` from source code.

    git clone https://github.com/zladovan/luster.git
    cd luster
    go install

>You need to have [git](https://git-scm.com/downloads) and [golang](https://golang.org/dl/) installed locally.

## How to scrape fans of your Facebook your page

To download list of users who like or follow your page (fans) run

    luster fans my-page

Where `my-page` is the name (string identifier) of your Facebook page.

You will be asked for your email address and password which can be used to login to Facebook.

Alternatively you can provide credentials as part of the command 

    luster -u my@email.com -p 12345 fans my-page

>Note that used account need to have assigned one of the [available roles](https://www.facebook.com/help/289207354498410) on your page.

If everything went fine you can expect output similar to following

    TIME,KIND,ID,NAME,LINK
    1581665652,like,111111111,John Doe,https://www.facebook.com/111111111
    1581663355,like,222222222,Kal Peterson,https://www.facebook.com/222222222
    1581661970,follow,333333333,Nikol Kus,https://www.facebook.com/333333333

### Limitations

It looks Facebook limits endpoint for getting information about "fans" with **maximum of 7k results** even when paging is used.
