const inquirer = require("inquirer");
const fs = require("fs");

const dbOptions = [
  {
    name: "Mongo DB (recommended)",
    short: "Mongo DB",
    value: "mongo",
    conn: "mongodb://localhost:27017"
  },
  { name: "MySQL", value: "sql-mysql", conn: "root:my-secret-pw@/test" },
  {
    name: "Postgres",
    value: "sql-postgres",
    conn:
      "postgres://postgres:mysecretpassword@localhost/postgres?sslmode=disable"
  }
];

exports.handleNew = configTemplate => cmd => {
  console.log(
    "This utility walks you through creating a config.yaml file for your space-cloud project."
  );
  console.log(
    "It only covers the most essential configurations and suggests sensible defaults.\n"
  );
  console.log("Press ^C at any time to quit.");
  inquirer
    .prompt([
      {
        type: "input",
        name: "id",
        message: "Project ID",
        default: answers => {
          const cwd = process.cwd();
          const keys = cwd.split("/");
          const currentFolderName = keys[keys.length - 1];
          return currentFolderName
            .toLowerCase()
            .split(" ")
            .join("-");
        }
      },
      {
        type: "list",
        name: "primaryDb",
        message: "Choose a main database",
        choices: dbOptions
      },
      {
        type: "input",
        name: "conn",
        message: answers =>
          `Connection string of ${
            dbOptions.find(db => db.value === answers.primaryDb).name
          }`,
        default: answers =>
          dbOptions.find(db => db.value === answers.primaryDb).conn
      }
    ])
    .then(answers => {
      fs.writeFileSync("./config.yaml", configTemplate(answers));
      console.log(
        "\nSuccess! Created a config.yaml file in the current directory"
      );
      console.log(
        "It consists of the details you entered just now along with all other possible configurations being commented out for you to play with.\n"
      );
      console.log("Next steps: ");
      console.log(
        `\n1] Read docs from https://spaceuptech.com/docs/ and edit config file as per your needs`
      );
      console.log("\n2] run space-cli deploy --local --config config.yaml");
    });
};
