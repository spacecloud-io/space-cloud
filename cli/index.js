#!/usr/bin/env node

const program = require("commander");
const fs = require("fs");
const handlebars = require("handlebars");

const { handleNew } = require("./handle");
const configTemplate = handlebars.compile(
  fs.readFileSync("./handlebars/config.handlebars", "utf8")
);

handlebars.registerHelper("ifEquals", function(arg1, arg2, options) {
  return arg1 === arg2 ? options.fn(this) : options.inverse(this);
});

program.usage("<command>").version("0.1.0", "-v, --version");

program
  .command("new")
  .description("Creates a new project config file")
  .action(handleNew(configTemplate));

program.parse(process.argv);
