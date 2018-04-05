const puppeteer = require('puppeteer');
const express = require('express')
const bodyParser = require('body-parser')

const app = express()
const port = 12458

app.use(bodyParser.urlencoded({     // to support URL-encoded bodies
        extended: true
    })
);

app.post('/', async (req, res) => {
    console.log(req.body.url)
    if (req.body.url !== null) {
        await getContent(req.body.url, res)
    }
})


app.listen(port, err => {
    if (err) {
        console.log(err)
    }
    else {
        console.log("server started on ", port)
    }
})

const tasks = Array()
var browser = null

const errHander = function (reason) {
    console.log(reason)
}

!function(){
    puppeteer.launch().then((brow) => {
       browser = brow
    }).catch(errHander)  
}()


/**
 * well designed to avoid memory leak
 * @author Jason<m.jason.liu@outlook.com> @CGQAQ
 * @version 1.1
 * @since 1.0
 * @param {string} url 
 * @param {*} response 
 * 
 * 2018.4.6 2.34 a.m.
 * 
 * Reusing page, it's not gonna change issue
 * need to goto about:Blank first
 * https://github.com/GoogleChrome/puppeteer/issues/1969
 */
async function getContent(url, response) {

    // const browser = await puppeteer.launch();
    // const page = await browser.newPage();

    // res = await page.goto(url)
    // if (!res.ok()) {
    //     return "";
    // }
    // return await res.text();

    if (browser === null) {
        tasks.push({
            url: url,
            response: response
        })
    }
    else{
        var page = await browser.newPage()
        while((task = tasks.pop())){
            await page.goto(task.url, {waitUntil: "networkidle0"})
            task.response.send(await page.content())
        }
    
        await page.goto(url, {waitUntil: "networkidle0"})
        response.end(await page.content())
        page.close()
    }
}


