import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../providers/card_notifier_provider.dart';
import '../../core/themes/colors.dart';
import '../../core/ui/custom/custom_app_bar.dart';
import '../../core/ui/custom/custom_button.dart';
import 'column_widget.dart';
import 'kanban_status.dart';

class KanbanScreen extends ConsumerWidget {
  const KanbanScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final cards = ref.watch(cardNotifierProvider).cards.toList();
    return Scaffold(
      backgroundColor: AppColors.transparent168,
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
          Expanded(
            child: SingleChildScrollView(
              scrollDirection: Axis.horizontal,
              child: Padding(
                padding: const EdgeInsets.symmetric(
                  horizontal: 10,
                  vertical: 20,
                ),
                child: Row(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children:
                      KanbanStatus.values
                          .map(
                            (col) => ColumnWidget(
                              columnId: col.id,
                              title: col.title,
                            ),
                          )
                          .expand(
                            (widget) => [widget, const SizedBox(width: 20)],
                          )
                          .toList()
                        ..removeLast(),
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}
