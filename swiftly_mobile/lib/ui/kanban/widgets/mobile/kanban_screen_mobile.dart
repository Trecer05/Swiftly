import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../domain/kanban/models/card_item.dart';
import '../../../../domain/kanban/models/priority.dart';
import '../../../../providers/card_notifier_provider.dart';
import '../../../../providers/current_user_provider.dart';
import '../../../../providers/selected_column_status_provider.dart';
import '../../../core/themes/colors.dart';
import '../../../core/ui/custom/custom_app_bar_mobile.dart';
import '../../../core/ui/custom/custom_button.dart';
import '../card_details_widgets/cart_details_widget.dart';
import '../kanban_status.dart';
import 'column_widget_mobile.dart';

class KanbanScreenMobile extends ConsumerWidget {
  const KanbanScreenMobile({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final cards = ref.watch(cardNotifierProvider).cards.toList();
    return Scaffold(
      backgroundColor: AppColors.transparent168,
      body: Stack(
        children: [
          Padding(
            padding: const EdgeInsets.only(left: 16, right: 16, top: 64),
            child: Column(
              children: [
                CustomAppBarMobile(
                  title: 'Задачи',
                  quantity: cards.length,
                  buttons: const [SegmentedControlWidgetMobile()],
                ),
                const SizedBox(height: 20),
                SingleChildScrollView(
                  scrollDirection: Axis.horizontal,
                  child: Row(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children:
                        KanbanStatus.values
                            .map(
                              (col) => ColumnStatusButton(
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
                const SizedBox(height: 20),
                const Expanded(
                  child: SingleChildScrollView(child: ColumnWidgetMobile()),
                ),
              ],
            ),
          ),
          Positioned(
            bottom: 10,
            left: 0,
            right: 0,
            child: Center(
              child: CustomButton(
                prefixIcon: Icons.add,
                text: 'Добавить',
                gradient: true,
                onTap: () => _handleCreateCard(context, ref),
              ),
            ),
          ),
        ],
      ),
    );
  }

  void _handleCreateCard(BuildContext context, WidgetRef ref) {
    final currentUser = ref.watch(currentUserProvider);
    final selectedColumnStatus = ref.watch(selectedColumnStatusProvider);
    if (currentUser == null) return;
    final newCard = CardItem.create(
      userId: currentUser.id,
      priority: Priority.low,
      columnId: selectedColumnStatus,
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
