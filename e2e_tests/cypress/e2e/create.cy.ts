describe('create forecast', () => {
    beforeEach(() => {
        cy.visit("/");
    });

    it('via the form on the main page succeeds', () => {
        // The purpose of randomNum is to allow running the test multiple times.
        const randomNum = Math.floor(Math.random() * 10000);
        const title = `Is this a forecast created on the main page? (${randomNum})`;
        const nextYear = (new Date().getFullYear() + 1).toString()
        const resolves = `0101${nextYear}1200PM`;

        cy.findByLabelText('Title *').type(title);
        // For some reason when selecting the 'resolves' field the last
        // element (AM/PM) is selected, so we need to move 5 times to the left.
        cy.findByLabelText('Resolves *').type('{leftArrow}'.repeat(5)+resolves);
        cy.findByLabelText('1. Outcome *').type('Yes');
        cy.findByLabelText('1. Probability *').type('95');
        cy.findByLabelText('2. Outcome *').type('No');
        cy.findByLabelText('2. Probability *').type('5');
        cy.findByLabelText('Reason *').type('Just a hunch.');
        cy.findByRole("button", {name: "Add Forecast"}).click();
        cy.findByText(/Saved ".*" with ID .*./);

        cy.findByLabelText('forecasts').within(() => {
            cy.findByText(title).parent("tr").as('row');
            cy.get('@row').within(() => {
                cy.findByText('UNRESOLVED').should('exist');
            });
        });
    });

});
