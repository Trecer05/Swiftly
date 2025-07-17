import 'package:flutter/material.dart';

import 'column_widget.dart';

class BoardScreen extends StatelessWidget {
  const BoardScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color.fromARGB(255, 9, 30, 114),
      body: SingleChildScrollView(
        scrollDirection: Axis.horizontal,
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 20),
          child: Column(
            children: [
              SizedBox(
                height: 50,
                child: FloatingActionButton(
                  onPressed: () {},
                  child: const Text('Добавить'),
                ),
              ),
              const Row(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  ColumnWidget(columnId: 'todo', title: 'To Do'),
                  SizedBox(width: 20),
                  ColumnWidget(columnId: 'progress', title: 'In Progress'),
                  SizedBox(width: 20),
                  ColumnWidget(columnId: 'done', title: 'Done'),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }
}
