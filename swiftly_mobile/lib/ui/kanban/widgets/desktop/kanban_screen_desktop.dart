import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:swiftly_mobile/ui/core/ui/custom/filter_state.dart';

import '../../../../domain/kanban/models/card_item.dart';
import '../../../../domain/kanban/models/priority.dart';
import '../../../../providers/card_notifier_provider.dart';
import '../../../../providers/current_user_provider.dart';
import '../../../core/themes/colors.dart';
import '../../../core/ui/custom/custom_app_bar_desktop.dart';
import '../../../core/ui/custom/custom_button.dart';
import '../card_details_widgets/cart_details_widget.dart';
import 'column_widget_desktop.dart';
import '../kanban_status.dart';

class KanbanScreenDesktop extends ConsumerWidget {
  const KanbanScreenDesktop({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    // final cards = ref.watch(cardNotifierProvider).cards.toList();
    final cards = ref.watch(filteredCardsProvider); // ✅ тут уже фильтрованные
    return Scaffold(
      backgroundColor: AppColors.transparent168,
      body: Column(
        children: [
          CustomAppBarDesktop(
            title: 'Задачи',
            quantity: cards.length,
            buttons: [
              CustomButton(
                prefixIcon: Icons.search,
                gradient: false,
                onTap: () {},
              ),
              const SegmentedControlWidgetDesktop(),
              CustomButton(
                prefixIcon: Icons.add,
                text: 'Добавить',
                gradient: true,
                onTap: () => _handleCreateCard(context, ref),
              ),
            ],
          ),
          Expanded(
            child: Padding(
              padding: const EdgeInsets.all(20),
              child: Row(
                crossAxisAlignment: CrossAxisAlignment.start,
                children:
                    KanbanStatus.values
                        .map(
                          (col) => Expanded(
                            child: ColumnWidgetDesktop(
                              columnId: col.id,
                              title: col.title,
                            ),
                          ),
                        )
                        .expand((widget) => [widget, const SizedBox(width: 20)])
                        .toList()
                      ..removeLast(),
              ),
            ),
          ),
        ],
      ),
    );
  }

  void _handleCreateCard(BuildContext context, WidgetRef ref) {
    final currentUser = ref.watch(currentUserProvider);
    if (currentUser == null) return;
    final newCard = CardItem.create(
      userId: currentUser.id,
      priority: Priority.low,
      columnId: 'todo',
    );
    ref.read(cardNotifierProvider.notifier).addCard(newCard);
    showDialog(
      context: context,
      builder:
          (context) => Dialog(
            insetPadding: const EdgeInsets.all(16),
            backgroundColor: Colors.transparent,
            child: CartDetailsWidget(card: newCard),
          ),
    );
  }
}
