const puppeteer = require('puppeteer');
const express = require('express')
const bodyParser = require('body-parser')

const app = express()
const port = 12458

app.use(bodyParser.urlencoded({     // to support URL-encoded bodies
    extended: true
  })
);

app.post('/', (req, res, err) =>{
    console.log(req.body.url)
    if(req.body.url !== null){
        (async function(){
            const content = await getContent(req.body.url, res);
            res.send(content);
            res.end();
        })();
    }
})



app.listen(port, err => {
    if(err){
        console.log(err)
    }
    else{
        console.log("server started on ", port)
    }
})


function getContent(url, response){

    // const browser = await puppeteer.launch();
    // const page = await browser.newPage();

    // res = await page.goto(url)
    // if (!res.ok()) {
    //     return "";
    // }
    // return await res.text();

    const errHander = function(reason){
        console.log(reason)
    }

    puppeteer.launch().then(browser => {
        browser.newPage().then(page => {
            page.goto(url).then(res => {
                if(res.ok()){
                    res.text().then(str => {
                        response.end(str)
                    }).catch(errHander)
                }
                else{
                    console.log('请求错误，状态码：', res.status())
                }
            }).catch(errHander)
        }).catch(errHander)
    }).catch(errHander)
}


