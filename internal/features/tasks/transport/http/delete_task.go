package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/Kor1992/todo/internal/core/logger"
	core_http_request "github.com/Kor1992/todo/internal/core/transport/http/request"
	core_http_response "github.com/Kor1992/todo/internal/core/transport/http/response"
)

// DeleteTask    godoc
// @Summary      Удаление задачи
// @Description  Удалить существующую в системе задачу по её ID
// @Tags         tasks
// @Param        id   path  int  true "ID удаляемой задачи" Format(uuid)
// @Success      204  "Успешное удаление задачи"
// @Failure      400  {object}  core_http_response.ErrorResponse "Bad request"
// @Failure      404  {object}  core_http_response.ErrorResponse "Task not found"
// @Failure      500  {object}  core_http_response.ErrorResponse "Internal server error"
// @Router       /tasks/{id} [delete]
func (h *TasksHTTPHandler) DeleteTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTEPResponseHandler(log, rw)

	taskId, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get task id path value")
		return
	}

	if err := h.tasksService.DeleteTask(ctx, taskId); err != nil {
		responseHandler.ErrorResponse(err, "failed to delete task")
		return
	}

	responseHandler.NoContentResponse()

}
