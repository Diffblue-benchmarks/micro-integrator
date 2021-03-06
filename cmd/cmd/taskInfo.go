/*
* Copyright (c) 2019, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
*
* WSO2 Inc. licenses this file to you under the Apache License,
* Version 2.0 (the "License"); you may not use this file except
* in compliance with the License.
* You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing,
* software distributed under the License is distributed on an
* "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
* KIND, either express or implied. See the License for the
* specific language governing permissions and limitations
* under the License.
 */

package cmd

import (
    "fmt"
    "github.com/olekukonko/tablewriter"
    "github.com/spf13/cobra"
    "github.com/wso2/micro-integrator/cmd/utils"
    "os"
)

var taskName string

// Show Task command related usage info
const showTaskCmdLiteral = "task"
const showTaskCmdShortDesc = "Get information about tasks"

var showTaskCmdLongDesc = "Get information about the Task specified by command line argument [task-name] If not specified, list all the tasks\n"

var showTaskCmdExamples = 
"Example:\n" +
"To get details about a specific task\n" +
"  " + programName + " " + showCmdLiteral + " " + showTaskCmdLiteral + " SampleTask\n\n" +
"To list all the tasks\n" +
"  " + programName + " " + showCmdLiteral + " " + showTaskCmdLiteral + "\n\n"

// taskShowCmd represents the task command
var taskShowCmd = &cobra.Command{
    Use:   showTaskCmdLiteral,
    Short: showTaskCmdShortDesc,
    Long:  showTaskCmdLongDesc + showTaskCmdExamples,
    Run: func(cmd *cobra.Command, args []string) {
        handleTaskCmdArguments(args)
    },
}

func init() {
    showCmd.AddCommand(taskShowCmd)
    taskShowCmd.SetHelpTemplate(showTaskCmdLongDesc + utils.GetCmdUsage(programName, showCmdLiteral, 
        showTaskCmdLiteral, "[task-name]") + showTaskCmdExamples + utils.GetCmdFlags("task(s)"))
}

func handleTaskCmdArguments(args []string) {
    utils.Logln(utils.LogPrefixInfo + "Show task called")
    if len(args) == 0 {
        executeListTasksCmd()
    } else if len(args) == 1 {
        if args[0] == "help" {
            printTaskHelp()
        } else {
            taskName = args[0]
            executeGetTaskCmd(taskName)
        }
    } else {
        fmt.Println("Too many arguments. See the usage below")
        printTaskHelp()
    }
}

func printTaskHelp() {
    fmt.Print(showTaskCmdLongDesc + utils.GetCmdUsage(programName, showCmdLiteral, showTaskCmdLiteral,
        "[task-name]") + showTaskCmdExamples + utils.GetCmdFlags("task(s)"))
}

func executeGetTaskCmd(taskname string) {

    finalUrl := utils.RESTAPIBase + utils.PrefixTasks + "?taskName=" + taskname

    resp, err := utils.UnmarshalData(finalUrl, &utils.Task{})

    if err == nil {
        // Printing the details of the Task
        task := resp.(*utils.Task)
        printTask(*task)
    } else {
        fmt.Println(utils.LogPrefixError+"Getting Information of the Task", err)
    }
}

// Print the details of a Task
// Name, Class, Group, Type and Trigger details
// @param task : Task object
func printTask(task utils.Task) {
    fmt.Println("Name - " + task.Name)
    fmt.Println("Trigger Type - " + task.Type)
    if task.Type == "cron" {
        fmt.Println("Cron Expression - " + task.TriggerCron)
    } else {
        fmt.Println("Trigger Count - " + task.TriggerCount)
        fmt.Println("Trigger Interval - " + task.TriggerInterval)
    }    
}

func executeListTasksCmd() {

    finalUrl := utils.RESTAPIBase + utils.PrefixTasks

    resp, err := utils.GetArtifactList(finalUrl, &utils.TaskList{})

    if err == nil {
        // Printing the list of available Tasks
        list := resp.(*utils.TaskList)
        printTaskList(*list)
    } else {
        utils.Logln(utils.LogPrefixError+"Getting List of Tasks", err)
    }
}

func printTaskList(taskList utils.TaskList) {

    if taskList.Count > 0 {
        table := tablewriter.NewWriter(os.Stdout)
        table.SetAlignment(tablewriter.ALIGN_LEFT)

        data := []string{"NAME", "TRIGGER TYPE", "COUNT", "INTERVAL", "CRON EXPRESSION"}
        table.Append(data)

        for _, task := range taskList.Tasks {
            data = []string{task.Name, task.Type, task.TriggerCount, task.TriggerInterval, task.TriggerCron}
            table.Append(data)
        }
        table.SetBorder(false)
        table.SetColumnSeparator("  ")
        table.Render()
    } else {
        fmt.Println("No Tasks found")
    }    
}