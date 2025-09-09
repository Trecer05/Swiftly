import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:swiftly_mobile/domain/kanban/models/priority.dart';
import 'package:swiftly_mobile/providers/card_notifier_provider.dart';
import 'package:swiftly_mobile/providers/current_user_provider.dart';

import '../../../core/themes/colors.dart';
import '../../../core/themes/theme.dart';
import '../../../../domain/kanban/models/card_item.dart';
import '../card_details_widgets/cart_details_widget.dart';
import '../../../core/ui/card_item/card_item_desktop.dart';

class ColumnWidgetDesktop extends ConsumerWidget {
  final String columnId;
  final String title;
  const ColumnWidgetDesktop({
    super.key,
    required this.columnId,
    required this.title,
  });

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final cards =
        ref
            .watch(cardNotifierProvider)
            .cards
            .where((card) => card.columnId == columnId)
            .toList();
    return SingleChildScrollView(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(title, style: AppTextStyles.style10),
          const SizedBox(height: 10),
          cards.isEmpty
              ? AddCardWidget(
                onCreate: () => _handleCreateCard(context, ref, columnId),
              )
              : Column(
                children: [
                  ...cards
                      .expand(
                        (card) => [
                          const SizedBox(height: 10),
                          CardItemDesktop(card: card),
                        ],
                      )
                      .skip(1),
                ],
              ),
        ],
      ),
    );
  }

  void _handleCreateCard(BuildContext context, WidgetRef ref, String columnId) {
    final currentUser = ref.watch(currentUserProvider);
    if (currentUser == null) return;
    final newCard = CardItem.create(
      userId: currentUser.id,
      priority: Priority.low,
      columnId: columnId,
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

class AddCardWidget extends StatefulWidget {
  final VoidCallback onCreate;
  const AddCardWidget({super.key, required this.onCreate});

  @override
  State<AddCardWidget> createState() => _AddCardWidgetState();
}

class _AddCardWidgetState extends State<AddCardWidget> {
  bool isHovered = false;

  @override
  Widget build(BuildContext context) {
    return MouseRegion(
      cursor: SystemMouseCursors.click,
      onEnter: (_) => setState(() => isHovered = true),
      onExit: (_) => setState(() => isHovered = false),
      child: GestureDetector(
        onTap: widget.onCreate,
        child: Container(
          width: double.infinity,
          height: 131,
          padding: const EdgeInsets.all(10),
          decoration: BoxDecoration(
            color: isHovered ? AppColors.white31 : AppColors.white15,
            borderRadius: BorderRadius.circular(15),
          ),
          child: const Center(
            child: Icon(Icons.add, color: AppColors.white, size: 18),
          ),
        ),
      ),
    );
  }
}
