import 'package:flutter/material.dart';

const List<Color> _defaultColors = [
  Colors.red,
  Colors.pink,
  Colors.purple,
  Colors.deepPurple,
  Colors.indigo,
  Colors.blue,
  Colors.lightBlue,
  Colors.cyan,
  Colors.teal,
  Colors.green,
  Colors.lightGreen,
  Colors.lime,
  Colors.yellow,
  Colors.amber,
  Colors.orange,
  Colors.deepOrange,
  Colors.brown,
  Colors.grey,
  Colors.blueGrey,
  Colors.black,
];

typedef PickerItem = Widget Function(Color color);

typedef PickerLayoutBuilder =
    Widget Function(BuildContext context, List<Color> colors, PickerItem child);

typedef PickerItemBuilder =
    Widget Function(
      Color color,
      bool isCurrentColor,
      void Function() changeColor,
    );

Widget _defaultLayoutBuilder(
  BuildContext context,
  List<Color> colors,
  PickerItem child,
) {
  return SizedBox(
    width: 300,
    height: 200,
    child: GridView.count(
      crossAxisCount: 6,
      crossAxisSpacing: 5,
      mainAxisSpacing: 5,
      children: [for (Color color in colors) child(color)],
    ),
  );
}

Widget _defaultItemBuilder(
  Color color,
  bool isCurrentColor,
  void Function() changeColor,
) {
  return Container(
    margin: const EdgeInsets.all(7),
    decoration: BoxDecoration(
      shape: BoxShape.circle,
      color: color,
      boxShadow: [
        BoxShadow(
          color: color.withValues(alpha: 0.8),
          offset: const Offset(1, 2),
          blurRadius: 5,
        ),
      ],
    ),
    child: Material(
      color: Colors.transparent,
      child: InkWell(
        onTap: changeColor,
        borderRadius: BorderRadius.circular(50),
        child: AnimatedOpacity(
          duration: const Duration(milliseconds: 210),
          opacity: isCurrentColor ? 1 : 0,
          child: const Icon(Icons.done, color: Colors.white),
        ),
      ),
    ),
  );
}

class CustomColorPicker extends StatefulWidget {
  final Color pickerColor;
  final ValueChanged<Color> onColorChanged;
  final List<Color> availableColors;
  final PickerLayoutBuilder layoutBuilder;
  final PickerItemBuilder itemBuilder;
  const CustomColorPicker({
    super.key,
    required this.pickerColor,
    required this.onColorChanged,
    this.availableColors = _defaultColors,
    this.layoutBuilder = _defaultLayoutBuilder,
    this.itemBuilder = _defaultItemBuilder,
  });

  @override
  State<CustomColorPicker> createState() => _CustomColorPickerState();
}

class _CustomColorPickerState extends State<CustomColorPicker> {
  Color? _currentColor;

  @override
  void initState() {
    _currentColor = widget.pickerColor;
    super.initState();
  }

  void changeColor(Color color) {
    setState(() => _currentColor = color);
    widget.onColorChanged(color);
  }

  @override
  Widget build(BuildContext context) {
    return widget.layoutBuilder(
      context,
      widget.availableColors,
      (Color color) => widget.itemBuilder(
        color,
        (_currentColor != null) ? (_currentColor?.value == color.value) : false,
        () => changeColor(color),
      ),
    );
  }
}
