import 'package:flutter/material.dart';

class AppBarCloud extends StatelessWidget {
  final String title;
  final String? currentDir;
  final int count;
  final VoidCallback onAdd;
  final Function(String)? onSearch;

  const AppBarCloud({
    required this.title,
    this.currentDir,
    required this.count,
    required this.onAdd,
    this.onSearch,
    super.key,
  });

  @override
Widget build(BuildContext context) {
  return Padding(
    padding: const EdgeInsets.only(top: 25, bottom: 0, left: 20, right: 24),
    child: Row(
      children: [
        Row(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              "Файлы",
              style: TextStyle(
                color: Colors.white,
                fontSize: 18,
                fontWeight: FontWeight.w700,
              ),
            ),
            Padding(
              padding: EdgeInsets.only(top: 0, left: 2),
              child: Text(
                '${count}',
                style: TextStyle(
                  color: Color(0xFF6DA8FF),
                  fontSize: 11,
                  fontWeight: FontWeight.w500,
                ),
              ),
            ),
          ],
        ),
        Expanded(
          child: Row(
            mainAxisAlignment: MainAxisAlignment.end,
            children: [
              CloudSearchField(onSearch: onSearch),
              SizedBox(width: 16),
              CloudSortButton(),
              SizedBox(width: 16),
              CloudAddButton(onAdd: onAdd),
            ],
          ),
        ),
      ],
    ),
  );
}
}

class CloudSearchField extends StatefulWidget {
  final Function(String)? onSearch;

  const CloudSearchField({this.onSearch, super.key});

  @override
  State<CloudSearchField> createState() => _CloudSearchFieldState();
}

class _CloudSearchFieldState extends State<CloudSearchField> {
  late TextEditingController _controller;
  late FocusNode _focusNode;
  bool _isExpanded = false;

  @override
  void initState() {
    super.initState();
    _controller = TextEditingController();
    _focusNode = FocusNode();

    _focusNode.addListener(_onFocusChange);
    _controller.addListener(_onSearchChange);
  }

  void _onFocusChange() {
    setState(() {
      if (!_focusNode.hasFocus && _controller.text.isEmpty) {
        _isExpanded = false;
      }
      if (_focusNode.hasFocus) {
        _isExpanded = true;
      }
    });
  }

  void _onSearchChange() {
    widget.onSearch?.call(_controller.text);
    
    if (_controller.text.isNotEmpty) {
      setState(() {
        _isExpanded = true;
      });
    }
  }

  void _clearSearch() {
    _controller.clear();
    widget.onSearch?.call('');
    _focusNode.unfocus();
    setState(() {
      _isExpanded = false;
    });
  }

  @override
  void dispose() {
    _controller.dispose();
    _focusNode.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return AnimatedContainer(
      duration: Duration(milliseconds: 300),
      width: _isExpanded ? 200 : 36, // ← УМЕНЬШЕНА ВЫСОТА
      height: 36, // ← ВЫСОТА ТАКАЯ ЖЕ КАК У ДРУГИХ КНОПОК
      decoration: BoxDecoration(
        color: Color(0x66182570),
        borderRadius: BorderRadius.circular(8),
        border: _isExpanded 
          ? Border.all(color: Color(0xFF3970F0), width: 1)
          : null,
      ),
      child: _isExpanded
          ? Row(
              children: [
                Expanded(
                  child: TextField(
                    controller: _controller,
                    focusNode: _focusNode,
                    style: TextStyle(color: Colors.white, fontSize: 12),
                    decoration: InputDecoration(
                      hintText: 'Поиск...',
                      hintStyle: TextStyle(color: Colors.grey[500], fontSize: 12),
                      border: InputBorder.none,
                      prefixIcon: Icon(Icons.search, color: Colors.white70, size: 18),
                      contentPadding: EdgeInsets.symmetric(vertical: 5, horizontal: 4),
                    ),
                  ),
                ),
                if (_controller.text.isNotEmpty)
                  Padding(
                    padding: const EdgeInsets.only(right: 4),
                    child: GestureDetector(
                      onTap: _clearSearch,
                      child: Icon(Icons.close, color: Colors.white70, size: 16),
                    ),
                  ),
              ],
            )
          : GestureDetector(
              onTap: () {
                setState(() {
                  _isExpanded = true;
                });
                _focusNode.requestFocus();
              },
              child: Center(
                child: Icon(Icons.search, color: Colors.white70, size: 20),
              ),
            ),
    );
  }
}

class CloudSortButton extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Container(
      height: 36, // ← ЕДИНАЯ ВЫСОТА
      padding: EdgeInsets.symmetric(horizontal: 12),
      decoration: BoxDecoration(
        color: Color(0x33182570),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(Icons.sort, color: Colors.white, size: 18),
          SizedBox(width: 6),
          Text(
            'По названию',
            style: TextStyle(color: Colors.white, fontSize: 12),
          ),
          SizedBox(width: 4),
          Icon(Icons.keyboard_arrow_down, color: Colors.white, size: 16),
        ],
      ),
    );
  }
}

class CloudAddButton extends StatelessWidget {
  final VoidCallback onAdd;

  const CloudAddButton({required this.onAdd, super.key});

  @override
  Widget build(BuildContext context) {
    return ElevatedButton(
      style: ElevatedButton.styleFrom(
        backgroundColor: Color(0xFF3970F0),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(8),
        ),
        padding: EdgeInsets.symmetric(horizontal: 12, vertical: 0),
        fixedSize: Size.fromHeight(36), // ← ЕДИНАЯ ВЫСОТА
      ),
      onPressed: onAdd,
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(Icons.add, color: Colors.white, size: 18),
          SizedBox(width: 4),
          Text(
            'Добавить',
            style: TextStyle(color: Colors.white, fontSize: 12),
          ),
        ],
      ),
    );
  }
}
