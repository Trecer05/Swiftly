import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../providers/card_notifier_provider.dart';
import '../../core/ui/custom/custom_app_bar.dart';
import '../../core/ui/custom/custom_button.dart';
import 'column_widget.dart';

class KanbanScreen extends ConsumerWidget {
  const KanbanScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final cards = ref.watch(cardNotifierProvider).cards.toList();
    return Scaffold(
      backgroundColor: const Color.fromARGB(255, 9, 30, 114),
      body: Column(
        children: [
          CustomAppBar(
            title: 'Задачи',
            quantity: cards.length,
            buttons: [
              CustomButton(
                prefixIcon: Icons.search,
                gradient: false,
                onTap: () {},
              ),
              CustomButton(
                prefixIcon: Icons.add,
                text: 'Добавить',
                gradient: true,
                onTap: () {},
              ),
            ],
          ),
          const Expanded(
            child: SingleChildScrollView(
              scrollDirection: Axis.horizontal,
              child: Padding(
                padding: EdgeInsets.symmetric(horizontal: 10, vertical: 20),
                child: Row(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    ColumnWidget(columnId: 'todo', title: 'To Do'),
                    SizedBox(width: 20),
                    ColumnWidget(columnId: 'progress', title: 'In Progress'),
                    SizedBox(width: 20),
                    ColumnWidget(columnId: 'done', title: 'Done'),
                  ],
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}
