const { By, Builder } = require('selenium-webdriver');
const { suite } = require('selenium-webdriver/testing');
const assert = require("assert");
const child_process = require('child_process');

const BASE_URL = 'http://localhost:8080';
const CLEOC_PATH = '../dist/cleoc_linux_amd64_v1/cleoc';

suite(function(env) {
    describe('Verify front page', function() {
        let driver;

        before(async function() {
            driver = await new Builder().forBrowser('firefox').build();

            // https://stackoverflow.com/a/32872753
            let child = child_process.spawnSync(
                CLEOC_PATH,
                [
                    'add',
                    'forecast',
                    '-t',
                    'Is this a test forecast?',
                    '-r',
                    '2022-12-01T15:00:00+01:00'
                ],
                { encoding : 'utf8' },
            );
            console.log("'cleoc add forecast' finished.");
            if(child.error) {
                console.log("ERROR: ",child.error);
            }
            console.log("stdout: ",child.stdout);
            console.log("stderr: ",child.stderr);
            console.log("exit code: ",child.status);
        });

        after(async () => await driver.quit());

        it('Verify title, header and new prediction appears', async function() {
            await driver.manage().setTimeouts({ implicit: 2000 });

            await driver.get(BASE_URL);

            let title = await driver.getTitle();

            assert.equal("Cleodora", title, "Title does not match");

            let header = await driver.findElement(By.css('h6'));

            let value = await header.getText();
            assert.equal("Cleodora", value, "Header does not match");

            let bodyText = await driver.findElement(By.tagName('body')).getText();
            console.log(bodyText);
            assert(
                bodyText.includes("Is this a test forecast?"),
                "Prediction created with cleoc not on the page",
            );
        });

    });
});
