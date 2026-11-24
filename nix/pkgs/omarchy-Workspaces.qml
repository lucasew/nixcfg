import QtQuick
import QtQuick.Layouts
import Quickshell
import Quickshell.Hyprland
import Quickshell.I3
import qs.Commons
import qs.Ui

BarWidget {
  id: root
  moduleName: "omarchy.workspaces"

  readonly property bool useI3: {
    var sway = Quickshell.env("SWAYSOCK")
    var i3 = Quickshell.env("I3SOCK")
    return (sway && sway.length > 0) || (i3 && i3.length > 0)
  }

  function workspaceEntries() {
    return root.useI3 ? I3.workspaces.values : Hyprland.workspaces.values
  }

  function workspaceNumber(workspace) {
    if (workspace === null)
      return -1
    return root.useI3 ? workspace.number : workspace.id
  }

  function workspaceByNumber(id) {
    var values = root.workspaceEntries()
    for (var i = 0; i < values.length; i++) {
      if (root.workspaceNumber(values[i]) === id)
        return values[i]
    }
    return null
  }

  function workspaceIds() {
    var ids = [1, 2, 3, 4, 5]
    var values = root.workspaceEntries()

    for (var i = 0; i < values.length; i++) {
      var id = root.workspaceNumber(values[i])
      if (id > 0 && id <= 10 && ids.indexOf(id) === -1)
        ids.push(id)
    }

    ids.sort(function (left, right) { return left - right })
    return ids
  }

  function focusedWorkspaceNumber() {
    if (root.useI3)
      return I3.focusedWorkspace !== null ? I3.focusedWorkspace.number : -1
    return Hyprland.focusedWorkspace !== null ? Hyprland.focusedWorkspace.id : -1
  }

  function workspaceOccupied(workspace) {
    if (workspace === null)
      return false
    if (root.useI3)
      return true
    return workspace.toplevels.values.length > 0
  }

  function focusWorkspace(id) {
    if (root.useI3) {
      I3.dispatch("workspace number " + id)
      return
    }
    if (!root.bar)
      return
    root.bar.run("hyprctl dispatch " + Util.shellQuote("hl.dsp.focus({ workspace = \"" + id + "\" })"))
  }

  readonly property real trailingGap: root.vertical ? 0 : Style.spaceReal(1.5)

  implicitWidth: grid.implicitWidth + trailingGap
  implicitHeight: grid.implicitHeight

  GridLayout {
    id: grid
    anchors.fill: parent
    anchors.rightMargin: root.trailingGap
    columns: root.vertical ? 1 : root.workspaceIds().length
    columnSpacing: root.vertical ? 0 : Style.space(1)
    rowSpacing: root.vertical ? Style.space(2) : 0

    Repeater {
      model: root.workspaceIds()

      WidgetButton {
        required property int modelData

        readonly property var workspace: root.workspaceByNumber(modelData)
        readonly property bool occupied: root.workspaceOccupied(workspace)
        readonly property bool focused: root.focusedWorkspaceNumber() === modelData

        bar: root.bar
        text: focused ? "\uDB85\uDCFB" : (modelData === 10 ? "0" : String(modelData))
        opacity: occupied || focused ? 1 : 0.5
        horizontalMargin: 6
        verticalPadding: 6
        fixedWidth: root.vertical ? root.barSize : Style.space(20)
        fixedHeight: root.barSize
        onPressed: function () { root.focusWorkspace(modelData) }
      }
    }
  }
}
