import React, { Component } from "react";
import axios from "axios";
import { Card, Header, Form, Input, Icon, Message, Button } from "semantic-ui-react";

let endpoint = "http://localhost:8080";

class ToDoList extends Component {
  constructor(props) {
    super(props);

    this.state = {
      task: "",
      items: [],
      notification: ""
    };
  }

  componentDidMount() {
    this.getTask();
  }

  onChange = event => {
    this.setState({
      [event.target.name]: event.target.value
    });
  };

  onSubmit = () => {
    let { task } = this.state;
    if (task) {
      axios
        .post(
          endpoint + "/api/task",
          {
            task
          },
          {
            headers: {
              "Content-Type": "application/x-www-form-urlencoded"
            }
          }
        )
        .then(res => {
          this.getTask();
          this.setState({
            task: "",
            notification: "task " + task + " has been created"
          });
          this.setState({...this.state, })
        });
    }
  };

  getTask = () => {
    axios.get(endpoint + "/api/task").then(res => {
      if (res.data.todo.tasks) {
        this.setState({
          items: res.data.todo.tasks.map((item) => {
            let color = "yellow";

            if (item.status) {
              color = "green";
            }
            return (
              <Card key={item.id} color={color} fluid>
                <Card.Content>
                  <Card.Header textAlign="left">
                    <div style={{ wordWrap: "break-word" }}>{item.task}</div>
                  </Card.Header>

                  <Card.Meta textAlign="right">
                    <Icon
                      name="check circle"
                      color="green"
                      onClick={() => this.updateTask(item.id)}
                    />
                    <span style={{ paddingRight: 10 }}>Done</span>
                    <Icon
                      name="undo"
                      color="yellow"
                      onClick={() => this.undoTask(item.id)}
                    />
                    <span style={{ paddingRight: 10 }}>Undo</span>
                    <Icon
                      name="delete"
                      color="red"
                      onClick={() => this.deleteTask(item.id)}
                    />
                    <span style={{ paddingRight: 10 }}>Delete</span>
                  </Card.Meta>
                </Card.Content>
              </Card>
            );
          })
        });
      } else {
        this.setState({
          items: []
        });
      }
    });
  };

  updateTask = id => {
    axios
      .put(endpoint + "/api/task/" + id, {
        headers: {
          "Content-Type": "application/x-www-form-urlencoded"
        }
      })
      .then(res => {
        this.getTask();
        this.setState({...this.state, notification: "task " + id + " has been updated"})
      });
  };

  undoTask = id => {
    axios
      .put(endpoint + "/api/undoTask/" + id, {
        headers: {
          "Content-Type": "application/x-www-form-urlencoded"
        }
      })
      .then(res => {
        this.getTask();
        this.setState({...this.state, notification: "task " + id + " has been updated"})
      });
  };

  deleteTask = id => {
    axios
      .delete(endpoint + "/api/deleteTask/" + id, {
        headers: {
          "Content-Type": "application/x-www-form-urlencoded"
        }
      })
      .then(res => {
        this.getTask();
        this.setState({...this.state, notification: "task " + id + " has been deleted"})
      });
  };

  deleteAllTask = () => {
    axios
      .delete(endpoint + "/api/deleteAllTask", {
        headers: {
          "Content-Type": "application/x-www-form-urlencoded"
        }
      })
      .then(res => {
        this.getTask();
        this.setState({...this.state, notification: "tasks has been purged"})
      });
  };

  deleteTask = id => {
    axios
      .delete(endpoint + "/api/deleteTask/" + id, {
        headers: {
          "Content-Type": "application/x-www-form-urlencoded"
        }
      })
      .then(res => {
        this.getTask();
        this.setState({...this.state, notification: "task " + id + " has been deleted"})
      });
  };

  removeNotification = () => {
    this.setState({...this.state, notification: ""})
  }

  render() {
    return (
      <div>
        <div className="row">
          <Header className="header" as="h2">
            TO DO LIST
          </Header>
        </div>
        {this.state.notification && <div className="row">
          <Message
            onDismiss={this.removeNotification}
            color='green'
            content={this.state.notification}
          />
        </div>
        }
        <div className="row">
          <Form onSubmit={this.onSubmit}>
            <Input
              type="text"
              name="task"
              onChange={this.onChange}
              value={this.state.task}
              fluid
              placeholder="Create Task"
            />
            {/* <Button >Create Task</Button> */}
          </Form>
        </div>
        <div className="row">
          <Card.Group>{this.state.items}</Card.Group>
        </div>
        <div className="row">
          <Button icon labelPosition='left' floated='right' color='red'
            onClick={() => this.deleteAllTask()}>
            <Icon
              bordered
              name="delete"
            />
            Delete all
          </Button>
        </div>
      </div>
    );
  }
}

export default ToDoList;


