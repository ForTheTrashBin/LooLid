***WORK IN PROGRESS***

# LooLid (A Directory-Based File Transformer)

Loolid is a program that allows you to copy one directory structure to another. Using various rules, you can control whether and how individual files or folders are copied. These rules determine whether individual files or directories are copied at all, or whether the contents of individual files should be modified before the copy process. The directory structure, directory names, and file names are always preserved.

## The basic configuration

The program offers very few options for direct configuration. Almost all settings necessary for the program to run correctly are fixed and cannot be customized. As a result, there are several requirements that must be met for the program to work.

- All settings are based on the current working directory.
- The root of the source and destination directories must be inside the working directory.
- The name of the source directory is 'content'.
- The name of the destination directory is 'target'
- The program is shipped under the name "LooLid" or "LooLid.exe" (on Windows)
- Possible configuration files are named 'LooLid.config'
- The user can rename the program file; the names of the configuration files are automatically adjusted accordingly
    - Example 1: LooLid(.exe) renamed to Donald.(exe). The configuration files must be named Donald.config
    - Example 2: LooLid(.exe) renamed to DAISY.(exe). The configuration files must be named DAISY.config
    
    - Example 3: LooLid(.exe) renamed to goofY.(exe). The configuration files must be named goofY.config

- Program features to look for
    - Upper- and lowercase letters are **strictly observed** (even on Windows).
    - The program will never ever make any changes to the source directory (content).
    
    - The program can change the destination directory (target) at any time — or even delete it.

## The configuration file(s)

Each directory or subdirectory may contain a configuration file. The baisc format of configuration file is TOML. Each configuration file can contain a frontmatter and up to two filter settings.

### Filter

- Filters may be specified in each configuration file.
- These filters apply to file and directory names in the current directory.
- These filters can be defined as blacklists or whitelists.
- Filters in the whitelist can remove filters in the blacklist.
- The filters in the black- and whitelists are inherited by the subdirectories.

- Inherited filters can be further expanded or modified in subdirectories.

### Frontmatter

- aaa
- bbb
- ccc

***WORK IN PROGRESS***
