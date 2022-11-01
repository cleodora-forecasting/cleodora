const { By, Builder } = require('selenium-webdriver');
const { suite } = require('selenium-webdriver/testing');
const assert = require("assert");

const BASE_URL = 'http://localhost:3000';

suite(function(env) {
    describe('Verify front page', function() {
        let driver;

        before(async function() {
            driver = await new Builder().forBrowser('firefox').build();
        });

        after(async () => await driver.quit());

        it('Verify title and header', async function() {
            await driver.get(BASE_URL);

            let title = await driver.getTitle();
            assert.equal("Cleodora", title);

            await driver.manage().setTimeouts({ implicit: 500 });

            let header = await driver.findElement(By.css('h2'));

            let value = await header.getText();
            assert.equal("Cleodora", value);
        });

    });
});
