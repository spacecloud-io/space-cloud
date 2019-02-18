#!/usr/bin/env node

const program = require("commander");
const inquirer = require("inquirer");
var emoji = require("node-emoji");
var colors = require("colors");
const boxen = require("boxen");

// function pad(txt, left = 0, right = 0) {
//   return (
//     Array(left)
//       .fill(" ")
//       .join("") +
//     txt +
//     Array(right)
//       .fill(" ")
//       .join("")
//   );
// }

program.usage("<command>").version("0.1.0", "-v, --version");

program
  .command("new")
  .usage("<project-name>")
  .description("Creates a new project config file")
  .option(
    "-q, --quick",
    "Quick mode skips default step by step guide and generates a boiler plate config file (Ideal if you know space cloud already!)"
  )
  .action(function(cmd) {
    if (typeof cmd !== "string") {
      console.log("Please provide a project name");
      return;
    }

    inquirer
      .prompt([
        {
          type: "list",
          name: "primaryDb",
          message: answers => {
            console.log(
              boxen(
                colors.cyan(
                  `${emoji.get("wave")} there! Welcome to space cloud!`
                ),
                { padding: 1, float: "center", borderColor: "", margin: 2 }
              )
            );

            console.log(
              `Let's build ${cmd} with Space Cloud ${emoji.get("rocket")}`
            );
            console.log(
              "Let us first select a primary database for your project"
            );
            console.log("What's a primary database?");
            console.log(
              "The crud module of space-cloud allows you to store data for your project in any database dynamically via the rest apis"
            );
            console.log(
              "However other modules (for eg user management) need a fixed database to store its data"
            );
            console.log("This is what we call as primary database");
            console.log(
              "You should not change a primary database once its selected"
            );
            console.log(
              "However you can always add other databases to the crud module to store data"
            );
            return "Select a primary database".padStart(10);
          },
          choices: [
            {
              name: "Mongo DB (recommended)",
              value: "mongo",
              short: "Mongo DB"
            },
            { name: "MySql", value: "mysql" },
            { name: "Postgres", value: "postgres" }
          ]
        },
        {
          type: "confirm",
          name: "isRealtimeEnabled",
          message: "Enable Realtime",
          default: true
        }
      ])
      .then(answers => {
        console.log("Answers", answers);
      });
  })
  .name("new <project-name>");

program
  .command("start")
  .description("Makes space cloud project up and running for use")
  .option("-c, --config", "Path ");

program.parse(process.argv);
