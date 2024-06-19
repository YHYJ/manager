/*
File: define_selector.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2024-04-10 13:33:59

Description: 定义选择器

- Init, Update, View 是不可或缺的方法
*/

package general

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gookit/color"
)

var quietKey = "q"                // 默认的退出键
var selectorType = "program name" // 选择器主题

// model 结构体，选择器的数据
type model struct {
	choices  []string         // 所有选项
	cursor   int              // 光标当前所在选项的索引
	selected map[int]struct{} // 已选中选项，key 为选项 choices 的索引。使用 map 便于判断指定选项是否已被选中
}

// initialModel 初始化 model
//
// 参数：
//   - choices: 可选项
//
// 返回：
//   - model
func initialModel(choices []string) model {
	allChoices := []string{color.Sprintf("%s%s", SelectAllFlag, FgLightYellowText(SelectAllTips))}
	allChoices = append(allChoices, choices...)

	return model{
		choices:  allChoices,
		cursor:   0,
		selected: make(map[int]struct{}),
	}
}

// Init model 结构体的初始化方法，是 BubbleTea 框架中的一个特殊方法
func (m model) Init() tea.Cmd {
	// 返回 nil 意味着不需要 I/O 操作
	return nil
}

// Update model 结构体的更新方法，是 BubbleTea 框架中的一个特殊方法
//
// 参数：
//   - msg: 包含来自 I/O 操作结果的数据，出发更新功能，并以此出触发 UI 绘制
//
// 返回：
//   - model: 更新后的 model
//   - tea.Cmd: 一个 I/O 操作，完成后会返回一条消息，如果为 nil 则被视为无操作
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// 监控按键事件
	case tea.KeyMsg:
		// 对按下的相应按键做出对应反应
		switch keyPress := msg.String(); keyPress {
		case quietKey, "ctrl+c", "esc":
			// 取消所有选中
			m.selected = make(map[int]struct{})
			return m, tea.Quit
		case "up", "k":
			// 光标向上移动，循环
			m.cursor--
			m.fixCursor()
		case "down", "j":
			// 光标向下移动，循环
			m.cursor++
			m.fixCursor()
		case " ":
			if m.cursor == 0 { // 选中“全选”项
				// 在全选和取消全选之间切换
				if len(m.selected) == len(m.choices)-1 { // 取消全选
					m.selected = make(map[int]struct{})
				} else { // 全选
					m.selected = make(map[int]struct{})
					for i := 1; i < len(m.choices); i++ {
						m.selected[i] = struct{}{} // 选中所有非 "Select All" 选项
					}
				}
			} else { // 选中其他选项
				// 判断当前光标所在选项是否被选中
				_, ok := m.selected[m.cursor]
				if ok { // 已选中，取消选中
					delete(m.selected, m.cursor)
				} else { // 未选中，选中
					m.selected[m.cursor] = struct{}{}
				}
			}
		case "enter":
			// 回车执行
			return m, tea.Quit
		}
	}
	// 将更新后的 model 返回给 BubbleTea 进行处理
	return m, nil
}

// View model 结构体的视图方法，是 BubbleTea 框架中的一个特殊方法
//
// 返回：
//   - string: 绘制内容
func (m model) View() string {
	// 构建显示内容
	s := strings.Builder{}
	s.WriteString(color.Sprintf("%s\n", strings.Repeat(Separator1st, len(MultiSelectTips)+len(selectorType))))
	s.WriteString(color.Sprintf(QuestionText(MultiSelectTips), selectorType))
	s.WriteString(color.Sprintf(SecondaryText(QuietTips), quietKey))
	s.WriteString(color.Sprintf("%s\n", strings.Repeat(Separator1st, len(MultiSelectTips)+len(selectorType))))

	// 对所有选项进行迭代
	SelectedFlag = SuccessText(SelectedFlag) // 选中状态标识
	selectedCount := 1                       // 已选中选项计数
	for i, choice := range m.choices {
		// 检查光标是否指向当前选项，默认未指向
		cursorFlag := CursorOffFlag // 未指向当前选项
		if m.cursor == i {
			cursorFlag = CursorOnFlag                          // 指向当前选项
			choice = color.Sprintf("\x1b[7m%s\x1b[0m", choice) // 光标所在选项着色
		}
		// 检查当前选项是否被选中
		checked := UnselectedFlag // 未选中
		if i == 0 {               // 如果当前选项索引为 0 即是 "Select All"
			if len(m.selected) == len(m.choices)-1 { // 所有选项都已选中（不包括 "Select All" 这个特殊选项）
				checked = SelectedFlag // 已选中
			}
		} else { // 如果当前选项索引不为 0 即非 "Select All"
			if _, ok := m.selected[i]; ok { // 当前选项的索引在已选中选项中
				checked = SelectedFlag                                                                           // 已选中
				choice = color.Sprintf("%s %s", FgLightBlueText(choice), SecondaryText("(", selectedCount, ")")) // 已选中选项着色
				selectedCount++
			}
		}
		s.WriteString(color.Sprintf("%s [%s] %s\n", cursorFlag, checked, choice))
	}
	s.WriteString(color.Sprintf("%s\n", strings.Repeat(Separator1st, len(MultiSelectTips)+len(selectorType))))
	return s.String()
}

// fixCursor 修正光标位置
func (m *model) fixCursor() {
	if m.cursor > len(m.choices)-1 {
		m.cursor = 0
	} else if m.cursor < 0 {
		m.cursor = len(m.choices) - 1
	}
}

// MultipleSelectionFilter 多选筛选器，接受一个可选项切片，返回一个已选项切片，允许全选
//
// 参数：
//   - choices: 可选项
//
// 返回：
//   - 已选项
//   - 错误信息
func MultipleSelectionFilter(choices []string) ([]string, error) {
	program := tea.NewProgram(initialModel(choices))

	// 返回 tea.Model
	initModel, err := program.Run()
	if err != nil {
		return nil, err
	}

	selectedChoices := []string{}
	// 将 tea.Model 断言为自定义 model
	if myModel, ok := initModel.(model); ok {
		if len(myModel.choices)-1 > 0 { // -1 除去 "Select All" 的干扰
			for index := range myModel.selected {
				selectedChoices = append(selectedChoices, myModel.choices[index])
			}
		}
	} else {
		return nil, fmt.Errorf("initModel is not model")
	}
	return selectedChoices, nil
}
