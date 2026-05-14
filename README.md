# SF-Helper CLI

A lightweight productivity tool written in Go to automate repetitive Salesforce CLI tasks.

## Features

- **External App Flow**: Automates the `package.xml` creation, retrieval from a source org, and deployment to a target org for External Client Applications.
- **Password Automation**: Generates an org password and automatically parses the output to copy the Login and Password directly to your macOS clipboard.

## Prerequisites

- [Go](https://go.dev/doc/install) (1.19+)
- [Salesforce CLI (sf)](https://developer.salesforce.com/tools/salesforcecli)
- macOS (uses `pbcopy` for clipboard integration)

## Installation

1. Clone the repository:
   ```bash
   git clone [https://github.com/yourusername/sf-helper.git](https://github.com/yourusername/sf-helper.git)
   cd sf-helper
